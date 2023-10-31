package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty"`
	Mobile     string             `json:"mobile,omitempty" bson:"mobile,omitempty"`
	Password   []byte             `json:"password,omitempty" bson:"password,omitempty"`
	Nickname   string             `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Avatar     string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Describe   string             `json:"describe,omitempty" bson:"describe,omitempty"`
	Status     string             `json:"status,omitempty" bson:"status,omitempty"` /* active / inactive */
	LastActive time.Time          `json:"lastActive,omitempty" bson:"lastActive,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
