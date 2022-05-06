package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Run struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	StartTime int64              `json:"startTime,omitempty" bson:"startTime,omitempty"`
	EndTime   int64              `json:"endTime,omitempty" bson:"endTime,omitempty"`
	SessionId primitive.ObjectID `json:"sessionId,omitempty" bson:"sessionId,omitempty"`
	ThingID   primitive.ObjectID `json:"thingId,omitempty" bson:"thingId,omitempty"`
}
