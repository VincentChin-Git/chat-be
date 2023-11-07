package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Msg struct {
	Id          *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderId    *primitive.ObjectID `json:"senderId,omitempty" bson:"senderId,omitempty"`
	ReceiveId   *primitive.ObjectID `json:"receiveId,omitempty" bson:"receiveId,omitempty"`
	Content     string              `json:"content,omitempty" bson:"content,omitempty"`
	ContentType string              `json:"contentType,omitempty" bson:"contentType,omitempty"` /* text / image / video */
	Status      string              `json:"status,omitempty" bson:"status,omitempty"`           /* sent / received / read / recalled / deletedS / deletedR / deletedAll */
	CreatedAt   *time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
