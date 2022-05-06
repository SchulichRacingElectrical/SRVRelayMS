package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartDate int64              `json:"startDate,omitempty" bson:"startDate,omitempty"`
	EndDate   int64              `json:"endDate,omitempty" bson:"endDate,omitempty"`
	ThingID   primitive.ObjectID `json:"thingId,omitempty" bson:"thingId,omitempty"`
}
