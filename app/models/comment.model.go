package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreationDate int64              `json:"creationDate,omitempty" bson:"creationDate,omitempty"`
	Content      string             `json:"content,omitempty" bson:"content,omitempty"`
	UserID       primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Type         string             `json:"type,omitempty" bson:"type,omitempty"`
	AssociatedId primitive.ObjectID `json:"associatedId,omitempty" bson:"associatedId,omitempty"`
}
