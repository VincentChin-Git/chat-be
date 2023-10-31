package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectId(s string) primitive.ObjectID {
	objId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		fmt.Println("Error convert objectId")
		return objId
	}
	return objId
}

func ToSkipRow(page int, pageSize int) int {
	return (page - 1) * pageSize
}
