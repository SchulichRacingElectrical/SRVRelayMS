package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreationDate time.Time          `json:"startDate,omitempty" bson:"startDate,omitempty"`
	Content      string             `json:"content,omitempty" bson:"content,omitempty"`
	UserID       primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
}
