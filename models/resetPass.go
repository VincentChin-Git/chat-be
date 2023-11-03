package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResetPass struct {
	Id         *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     *primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	VerifyCode string              `json:"verifyCode,omitempty" bson:"verifyCode,omitempty"`
	Status     string              `json:"status,omitempty" bson:"status,omitempty"` /* pending / completed */
	CreatedAt  *time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt  *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
