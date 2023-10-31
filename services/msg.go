package services

import (
	"chat-be/models"
	"chat-be/storage"
	"chat-be/utils"
	"context"
	"errors"
	"fmt"

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
