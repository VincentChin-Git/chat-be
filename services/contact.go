package services

import (
	"chat-be/models"
	"chat-be/storage"
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
	ContactId  primitive.ObjectID `json:"contactId,omitempty" bson:"contactId,omitempty"`
	Mobile     string             `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Nickname   string             `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Avatar     string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Describe   string             `json:"describe,omitempty" bson:"describe,omitempty"`
	LastActive time.Time          `json:"lastActive,omitempty" bson:"lastActive,omitempty"`
}

func GetContact(_id string, skip int, limit int) ([]ContactElemRes, error) {
	contactDoc := storage.ClientDatabase.Collection("contacts")

	// get contact list and populate each row with corresponding user info
	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: []bson.M{
				{"userId": _id},
				{"status": "active"},
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
	cur, err := contactDoc.Aggregate(context.Background(), mongo.Pipeline{matchStage, skipStage, limitStage, lookupStage, unwindStage})
	if err != nil {
		fmt.Println("errGetContactList", err.Error())
		return nil, errors.New("")
	}

	defer cur.Close(context.Background())

	var contactList []ContactElemRes
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

		res := ContactElemRes{
			ContactId:  item.ContactId,
			Mobile:     item.ContactInfo.Mobile,
			Avatar:     item.ContactInfo.Avatar,
			Nickname:   item.ContactInfo.Nickname,
			Describe:   item.ContactInfo.Describe,
			LastActive: item.ContactInfo.LastActive,
		}

		contactList = append(contactList, res)
	}
	return contactList, nil
}
