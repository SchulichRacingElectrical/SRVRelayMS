package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Run struct {
	ID					primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
	Name				string								`json:"name,omitempty" bson:"name,omitempty"`
	StartTime		string								`json:"startTime,omitempty" bson:"startTime,omitempty"`
	EndTime			string								`json:"endTime,omitempty" bson:"endTime,omitempty"`
	CommentIDs	[]primitive.ObjectID	`json:"commentIds,omitempty" bson:"commentIds,omitempty"`
}