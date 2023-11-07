package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	Id            *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId        *primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	ContactId     *primitive.ObjectID `json:"contactId,omitempty" bson:"contactId,omitempty"`
	RelativePoint int                 `json:"relativePoint,omitempty" bson:"relativePoint,omitempty"`
	Status        string              `json:"status,omitempty" bson:"status,omitempty"` /* active / inactive */
	CreatedAt     *time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
