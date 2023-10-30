package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type test struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`             /* _id */
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"` /* createdAt */
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"` /* updatedAt */
}
