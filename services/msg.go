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

func GetMsgs(userId string, contactId string, skip int, limit int) ([]models.Msg, error) {
	if userId == "" || contactId == "" {
		return []models.Msg{}, errors.New("Invalid Info")
	}
	msgDoc := storage.ClientDatabase.Collection("msgs")

	uOId := utils.ToObjectId(userId)
	cOId := utils.ToObjectId(contactId)

	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"$or": []bson.M{
					{"senderId": uOId, "receiveId": cOId, "status": bson.M{
						"$not": bson.M{
							"$in": []string{"recalled", "deletedS", "deletedAll"},
						},
					}},
					{"receiveId": uOId, "senderId": cOId, "status": bson.M{
						"$not": bson.M{
							"$in": []string{"recalled", "deletedR", "deletedAll"},
						},
					}},
				},
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
	cur, err := msgDoc.Aggregate(context.Background(), mongo.Pipeline{matchStage, skipStage, limitStage})

	if err != nil {
		fmt.Println(err, "errGetMsg")
		return []models.Msg{}, errors.New("")
	}

	defer cur.Close(context.Background())

	var res []models.Msg
	for cur.Next(context.Background()) {
		var item models.Msg
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println(err, "errGetMsg1")
			return []models.Msg{}, errors.New("")
		}
		res = append(res, item)
	}

	return res, nil
}

func UpdateMsgStatus(_id string, status string, userId string) error {
	if _id == "" || status == "" {
		return errors.New("Invalid Info")
	}
	msgDoc := storage.ClientDatabase.Collection("msgs")

	validStatus := []string{"received", "read", "recalled", "deletedS", "deletedR"} // deletedAll
	isErr := true
	for _, char := range validStatus {
		if char == status {
			isErr = false
			break
		}
	}
	if isErr {
		return errors.New("Invalid Info")
	}

	// validate update status
	var msgData models.Msg
	cur := msgDoc.FindOne(context.Background(), bson.M{"_id": utils.ToObjectId(_id)})
	if cur.Err() != nil {
		fmt.Println(cur.Err().Error())
		return errors.New("")
	}

	err := cur.Decode(&msgData)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("")
	}

	statusErr := true
	switch status {
	case "received":
		if msgData.Status == "sent" {
			statusErr = false
		}
	case "read":
		if msgData.Status == "received" {
			statusErr = false
		}
	case "recalled":
		if msgData.Status == "sent" || msgData.Status == "received" || msgData.Status == "read" || msgData.Status == "deletedR" {
			statusErr = false
		}
	case "deletedS":
		if msgData.Status != "deletedS" && userId == (*msgData.SenderId).Hex() {
			statusErr = false
			if msgData.Status == "deletedR" {
				status = "deletedAll"
			}
		}
	case "deletedR":
		if msgData.Status != "deletedR" && userId == (*msgData.ReceiveId).Hex() {
			statusErr = false
			if msgData.Status == "deletedS" {
				status = "deletedAll"
			}
		}
	default:
		break
	}

	if statusErr {
		return errors.New("Status Error")
	}

	_, err = msgDoc.UpdateByID(context.Background(), utils.ToObjectId(_id), bson.M{
		"$set": bson.M{
			"status": status,
		},
	})

	if err != nil {
		return errors.New("")
	}
	return nil
}

func SendMsg(senderId string, receiveId string, content string, contentType string) (string, error) {
	if senderId == "" || receiveId == "" || content == "" || contentType == "" {
		return "", errors.New("Invalid Info")
	}

	allowedContentType := []string{"text", "image", "video"}
	isTypeErr := true

	for _, t := range allowedContentType {
		if t == contentType {
			isTypeErr = false
			break
		}
	}

	if isTypeErr {
		fmt.Println("typeError", contentType)
		return "", errors.New("Invalid Message Type")
	}

	userDoc := storage.ClientDatabase.Collection("users")
	senderCount, errS := userDoc.CountDocuments(context.Background(), bson.M{"_id": utils.ToObjectId(senderId)})
	receiveCount, errR := userDoc.CountDocuments(context.Background(), bson.M{"_id": utils.ToObjectId(receiveId)})

	if errS != nil || errR != nil {
		return "", errors.New("")
	}
	if !(senderCount > 0 && receiveCount > 0) {
		return "", errors.New("Invalid User")
	}

	var msgData models.Msg

	senderOId := utils.ToObjectId(senderId)
	receiveOId := utils.ToObjectId(receiveId)
	timeNow := time.Now()
	msgData.SenderId = &senderOId
	msgData.ReceiveId = &receiveOId
	msgData.Content = content
	msgData.ContentType = contentType
	msgData.Status = "sent"
	msgData.CreatedAt = &timeNow
	msgData.UpdatedAt = &timeNow

	msgDoc := storage.ClientDatabase.Collection("msgs")
	result, err := msgDoc.InsertOne(context.Background(), msgData)
	if err != nil {
		return "", errors.New("")
	}

	_id := result.InsertedID.(primitive.ObjectID).Hex()

	return _id, nil
}

type GetOverviewMsgType struct {
	MsgData     models.Msg  `json:"msgData,omitempty"`
	ContactData models.User `json:"contactData,omitempty"`
}

func GetOverviewMsg(userId string, skip int, limit int) ([]GetOverviewMsgType, error) {
	if userId == "" {
		return []GetOverviewMsgType{}, errors.New("Invalid Info")
	}

	msgDoc := storage.ClientDatabase.Collection("msgs")

	// statusMatch := []string{"sent", "received", "read", ""}
	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"$or": []bson.M{
					{
						"senderId": utils.ToObjectId(userId),
						"status": bson.M{
							"$not": bson.M{"$in": []string{"recalled", "deletedS", "deletedAll"}},
						},
					},
					{
						"receiveId": utils.ToObjectId(userId),
						"status": bson.M{
							"$not": bson.M{"$in": []string{"recalled", "deletedR", "deletedAll"}},
						},
					},
				},
			},
		},
	}
	sortStage := bson.D{
		primitive.E{
			Key: "$sort", Value: bson.M{
				"createdAt": -1,
			},
		},
	}
	addFieldStage := bson.D{
		primitive.E{
			Key: "$addFields", Value: bson.M{
				"dependId": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": []string{"$senderId", userId}},
						"$receiveId",
						"$senderId",
					},
				},
			},
		},
	}
	groupStage := bson.D{
		primitive.E{
			Key: "$group", Value: bson.M{
				"_id":     "$dependId",
				"msgData": bson.M{"$first": "$$ROOT"},
			},
		},
	}
	lookupStage := bson.D{
		primitive.E{
			Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "users"},
				primitive.E{Key: "localField", Value: "_id"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "contactData"},
			},
		},
	}
	unwindStage := bson.D{
		primitive.E{
			Key: "$unwind", Value: bson.D{
				primitive.E{Key: "path", Value: "$contactData"},
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
	cur, err := msgDoc.Aggregate(context.Background(), mongo.Pipeline{
		matchStage,
		sortStage,
		addFieldStage,
		groupStage,
		lookupStage,
		unwindStage,
		skipStage,
		limitStage,
	})

	if err != nil {
		fmt.Println(err)
		return []GetOverviewMsgType{}, errors.New("")
	}

	var res []GetOverviewMsgType
	for cur.Next(context.Background()) {
		var temp, item GetOverviewMsgType
		err = cur.Decode(&temp)
		if err != nil {
			fmt.Println(err.Error(), "errGetMsgOverview")
			continue
		}

		item.MsgData = temp.MsgData
		item.ContactData.Id = temp.ContactData.Id
		item.ContactData.Mobile = temp.ContactData.Mobile
		item.ContactData.Nickname = temp.ContactData.Nickname
		item.ContactData.Avatar = temp.ContactData.Avatar
		item.ContactData.Describe = temp.ContactData.Describe
		item.ContactData.LastActive = temp.ContactData.LastActive

		res = append(res, item)
	}

	return res, nil
}

type GetUnreadMsgType struct {
	ContactId primitive.ObjectID `json:"contactId,omitempty" bson:"contactId,omitempty"`
	Unread    int                `json:"unread,omitempty" bson:"unread,omitempty"`
}

func GetUnreadMsg(userId string) ([]GetUnreadMsgType, error) {
	msgDoc := storage.ClientDatabase.Collection("msgs")
	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: bson.M{
				"receiveId": utils.ToObjectId(userId),
				"status": bson.M{
					"$in": []string{"sent", "received"},
				},
			},
		},
	}
	groupStage := bson.D{
		primitive.E{
			Key: "$group", Value: bson.M{
				"_id":       "$senderId",
				"contactId": bson.M{"$first": "$senderId"},
				"unread": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	unread := []GetUnreadMsgType{}
	cur, err := msgDoc.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	for cur.Next(context.Background()) {

		var tempUnread GetUnreadMsgType
		err := cur.Decode(&tempUnread)
		if err != nil {
			fmt.Println(err)
			return []GetUnreadMsgType{}, errors.New("")
		}
		unread = append(unread, tempUnread)
	}
	if err != nil {
		fmt.Println(err.Error())
		return []GetUnreadMsgType{}, errors.New("")
	}
	return unread, nil
}
