package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Run struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	StartDate time.Time            `json:"startDate,omitempty" bson:"startDate,omitempty"`
	EndDate   time.Time            `json:"endDate,omitempty" bson:"endDate,omitempty"`
	Comments  []primitive.ObjectID `json:"comments,omitempty" bson:"comments,omitempty"`
	ThingID   primitive.ObjectID   `json:"thingId,omitempty" bson:"thingId,omitempty"`
}
