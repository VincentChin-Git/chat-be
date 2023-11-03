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
			Key: "$or", Value: bson.A{
				[]bson.M{{"senderId": uOId}, {"receiveId": cOId}},
				[]bson.M{{"receiveId": uOId}, {"senderId": cOId}},
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
			return []models.Msg{}, errors.New("")
		}
		res = append(res, item)
	}

	return res, nil
}

func UpdateMsgStatus(_id string, status string) error {
	if _id == "" || status == "" {
		return errors.New("Invalid Info")
	}

	validStatus := []string{"failed", "sending", "sent", "received", "read", "recalled", "deletedS", "deletedR", "deletedAll"}
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

	msgDoc := storage.ClientDatabase.Collection("msgs")
	_, err := msgDoc.UpdateByID(context.Background(), utils.ToObjectId(_id), bson.M{"status": status})

	if err != nil {
		return errors.New("")
	}
	return nil
}

func SendMsg(senderId string, receiveId string, data models.Msg) (string, error) {
	if senderId == "" || receiveId == "" || data.Content == "" || data.ContentType == "" {
		return "", errors.New("Invalid Info")
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

	senderOId := utils.ToObjectId(senderId)
	receiveOId := utils.ToObjectId(receiveId)
	timeNow := time.Now()
	data.SenderId = &senderOId
	data.ReceiveId = &receiveOId
	data.Status = "sending"
	data.CreatedAt = &timeNow
	data.UpdatedAt = &timeNow

	msgDoc := storage.ClientDatabase.Collection("msgs")
	result, err := msgDoc.InsertOne(context.Background(), data)
	if err != nil {
		return "", errors.New("")
	}

	_id := result.InsertedID.(string)

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

	statusMatch := []string{""}
	matchStage := bson.D{
		primitive.E{
			Key: "$match", Value: []bson.M{
				{"$or": bson.A{
					bson.M{
						"senderId": utils.ToObjectId(userId),
						"status": bson.M{
							"$not": "deletedS",
						},
					},
					bson.M{
						"receiveId": utils.ToObjectId(userId),
						"status": bson.M{
							"$not": "deletedR",
						},
					},
				}},
				{"status": bson.M{"$in": statusMatch}},
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
		item.ContactData.Mobile = temp.ContactData.Mobile
		item.ContactData.Nickname = temp.ContactData.Nickname
		item.ContactData.Avatar = temp.ContactData.Avatar
		item.ContactData.Describe = temp.ContactData.Describe
		item.ContactData.LastActive = temp.ContactData.LastActive

		res = append(res, item)
	}

	return res, nil

}
