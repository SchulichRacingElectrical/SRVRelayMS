package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID						primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
	Name					string								`json:"name,omitempty" bson:"name,omitempty"`
	Description		string								`json:"description,omitempty" bson:"description,omitempty"`
	RunIds				[]primitive.ObjectID	`json:"runIds,omitempty" bson:"runIds,omitempty"`
	CommentIDs		[]primitive.ObjectID	`json:"commentIds,omitempty" bson:"commentIds,omitempty"`
}