package models

import "time"

type User struct {
	Username  string    `json:"username" bson:"username"`
	Mobile    string    `json:"mobile" bson:"mobile"`
	Password  []byte    `json:"password" bson:"password"`
	Nickname  string    `json:"nickname" bson:"nickname"`
	Avatar    string    `json:"avatar,omitempty" bson:"avatar,omitempty"` /* active / inactive */
	Describe  string    `json:"describe,omitempty" bson:"describe,omitempty"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
