package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThingOperator struct {
	ID						primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
	OperatorId		primitive.ObjectID		`json:"operatorId,omitempty" bson:"operatorId,omitempty"`
	ThingId				primitive.ObjectID		`json:"thingId,omitempty" bson:"thingId,omitempty"`
}