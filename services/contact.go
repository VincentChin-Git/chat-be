package services

import (
	"chat-be/models"
	"chat-be/storage"
	"chat-be/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type contactElem struct {
	ContactId   primitive.ObjectID `json:"contactId,omitempty" bson:"contactId,omitempty"`
	ContactInfo models.User        `json:"contactInfo,omitempty" bson:"contactInfo,omitempty"`
}

type ContactElemRes struct {
	ContactId  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Mobile     string             `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Nickname   string             `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Avatar     string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Describe   string             `json:"describe,omitempty" bson:"describe,omitempty"`
	LastActive time.Time          `json:"lastActive,omitempty" bson:"lastActive,omitempty"`
}

func GetContacts(_id string, sortBy string, skip int, limit int) ([]ContactElemRes, error) {
	contactDoc := storage.ClientDatabase.Collection("contacts")

	var sortKey string
	var sortValue int

	if sortBy == "name" {
		sortKey = "contactInfo.nickname"
		sortValue = 1
	} else if sortBy == "point" {
		sortKey = "relativePoint"
		sortValue = -1
	}

	// get contact list and populate each row with corresponding user info
	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"userId": utils.ToObjectId(_id),
				"status": "active",
			},
		},
	}
	sortStage := bson.D{
		primitive.E{
			Key: "$sort", Value: bson.M{
				sortKey: sortValue,
			},
		},
	}
	lookupStage := bson.D{
		primitive.E{
			Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "users"},
				primitive.E{Key: "localField", Value: "contactId"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "contactInfo"},
			},
		},
	}
	unwindStage := bson.D{
		primitive.E{
			Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$contactInfo"},
				primitive.E{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		},
	}
	skipStage := bson.D{
		primitive.E{
			Key: "$skip", Value: skip,
		},
	}
	limitStage := bson.D{
		primitive.E{
			Key: "$limit", Value: limit,
		},
	}
	cur, err := contactDoc.Aggregate(context.Background(), mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		sortStage,
		skipStage,
		limitStage,
	})
	if err != nil {
		fmt.Println("errGetContactList", err.Error())
		return nil, errors.New("")
	}

	defer cur.Close(context.Background())

	contactList := []ContactElemRes{}
	for cur.Next(context.Background()) {
		var item contactElem
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println("errGetContactElem", err.Error())
		}

		// get only active user
		if item.ContactInfo.Status != "active" {
			continue
		}

		lastActive := item.ContactInfo.LastActive
		res := ContactElemRes{
			ContactId:  item.ContactId,
			Mobile:     item.ContactInfo.Mobile,
			Avatar:     item.ContactInfo.Avatar,
			Nickname:   item.ContactInfo.Nickname,
			Describe:   item.ContactInfo.Describe,
			LastActive: *lastActive,
		}

		contactList = append(contactList, res)
	}
	return contactList, nil
}

func GetContact(param string, userId string, skip int, limit int) ([]ContactElemRes, error) {

	if param == "" {
		return []ContactElemRes{}, errors.New("Invalid Info")
	}

	contactDoc := storage.ClientDatabase.Collection("contacts")

	matchStage1 := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"userId": utils.ToObjectId(userId),
				"status": "active",
			},
		},
	}
	lookupStage := bson.D{
		primitive.E{
			Key: "$lookup", Value: bson.M{
				"from":         "users",
				"localField":   "contactId",
				"foreignField": "_id",
				"as":           "contactInfo",
			},
		},
	}
	unwindStage := bson.D{
		primitive.E{
			Key: "$unwind", Value: bson.M{
				"path":                       "$contactInfo",
				"preserveNullAndEmptyArrays": false,
			},
		},
	}
	matchStage2 := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"$or": []bson.M{
					{"contactInfo.nickname": bson.M{
						"$regex":   param,
						"$options": "i",
					}},
					{"contactInfo.mobile": bson.M{
						"$regex":   param,
						"$options": "i",
					}},
				},
			},
		},
	}
	sortStage := bson.D{
		primitive.E{
			Key: "$sort", Value: bson.M{
				"relativePoint": -1,
			},
		},
	}
	skipStage := bson.D{
		primitive.E{
			Key: "$skip", Value: skip,
		},
	}
	limitStage := bson.D{
		primitive.E{
			Key: "$limit", Value: limit,
		},
	}

	cur, err := contactDoc.Aggregate(context.Background(), mongo.Pipeline{
		matchStage1,
		lookupStage,
		unwindStage,
		matchStage2,
		sortStage,
		skipStage,
		limitStage,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	contactList := []ContactElemRes{}
	for cur.Next(context.Background()) {
		var item contactElem
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println("errGetContactElem", err.Error())
		}

		// get only active user
		if item.ContactInfo.Status != "active" {
			continue
		}

		lastActive := item.ContactInfo.LastActive
		res := ContactElemRes{
			ContactId:  item.ContactId,
			Mobile:     item.ContactInfo.Mobile,
			Avatar:     item.ContactInfo.Avatar,
			Nickname:   item.ContactInfo.Nickname,
			Describe:   item.ContactInfo.Describe,
			LastActive: *lastActive,
		}

		contactList = append(contactList, res)
	}
	return contactList, nil

}

func AddContact(userId string, contactId string) (models.Contact, error) {
	userDoc := storage.ClientDatabase.Collection("users")
	contactDoc := storage.ClientDatabase.Collection("contacts")

	// check user exist
	userList := []primitive.ObjectID{utils.ToObjectId(userId), utils.ToObjectId(contactId)}
	count, err := userDoc.CountDocuments(context.Background(), bson.M{"_id": bson.M{"$in": userList}})

	if err != nil {
		fmt.Println(err, "errFindUser")
		return models.Contact{}, errors.New("")
	}

	if count < 2 {
		return models.Contact{}, errors.New("Invalid User")
	}

	// check contact exist
	count, err = contactDoc.CountDocuments(context.Background(), bson.M{
		"userId":    utils.ToObjectId(userId),
		"contactId": utils.ToObjectId(contactId),
		"status":    "active",
	})

	if err != nil {
		fmt.Println(err, "existContact")
		return models.Contact{}, errors.New("")
	}

	if count != 0 {
		return models.Contact{}, errors.New("Contact Already Added")
	}

	userOId := utils.ToObjectId(userId)
	contactOId := utils.ToObjectId(contactId)
	timeNow := time.Now()
	newContact := models.Contact{
		UserId:        &userOId,
		ContactId:     &contactOId,
		RelativePoint: 1,
		Status:        "active",
		CreatedAt:     &timeNow,
		UpdatedAt:     &timeNow,
	}

	result, err := contactDoc.InsertOne(context.Background(), newContact)
	if err != nil {
		fmt.Println(err, "errAddContact", result.InsertedID)
		return models.Contact{}, errors.New("")
	}

	insertedId := result.InsertedID.(primitive.ObjectID)
	newContact.Id = &insertedId

	return newContact, nil
}

func RemoveContact(userId string, contactId string) error {

	contactDoc := storage.ClientDatabase.Collection("contacts")

	_, err := contactDoc.UpdateOne(context.Background(), bson.M{
		"userId":    utils.ToObjectId(userId),
		"contactId": utils.ToObjectId(contactId),
	}, bson.M{
		"$set": bson.M{
			"status":    "inactive",
			"updatedAt": time.Now(),
		},
	})

	if err != nil {
		fmt.Println(err, "errRemoveContact")
		return errors.New("")
	}

	return nil
}

func updatePoint(userId string, contactId string, isAdd bool, point int) error {

	contactDoc := storage.ClientDatabase.Collection("contacts")
	if point == 0 {
		point = 1
	}

	cur1 := contactDoc.FindOne(context.Background(), bson.M{
		"userId":    userId,
		"contactId": contactId,
	})

	cur2 := contactDoc.FindOne(context.Background(), bson.M{
		"userId":    contactId,
		"contactId": userId,
	})

	if cur1.Err() != nil || cur2.Err() != nil {
		fmt.Println(cur1.Err().Error(), cur2.Err().Error())
		return errors.New("")
	}

	var contact1, contact2 models.Contact

	err1 := cur1.Decode(&contact1)
	err2 := cur1.Decode(&contact2)

	if err1 != nil {
		fmt.Println(err1, "err1")
		return errors.New("")
	}

	var newPoint int
	if isAdd {
		newPoint = contact1.RelativePoint + point
	} else {
		newPoint = contact1.RelativePoint - point
	}
	_, err := contactDoc.UpdateByID(context.Background(), utils.ToObjectId(contact1.Id.Hex()), bson.M{
		"$set": bson.M{"relativePoint": newPoint},
	})
	if err != nil {
		fmt.Println(err, "err1")
		return errors.New("")
	}

	// maybe blank, so don't trigger error
	fmt.Println(err2, "err2")

	if err2 == nil {
		if isAdd {
			newPoint = contact2.RelativePoint + point
		} else {
			newPoint = contact2.RelativePoint - point
		}

		_, err := contactDoc.UpdateByID(context.Background(), utils.ToObjectId(contact2.Id.Hex()), bson.M{
			"$set": bson.M{"relativePoint": newPoint},
		})

		if err != nil {
			fmt.Println(err, "err2")
			return errors.New("")
		}
	}

	return nil
}
