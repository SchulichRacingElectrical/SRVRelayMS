package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChartPreset struct {
	ID						primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
	Name					string								`json:"name,omitempty" bson:"name,omitempty"`
	ThingId 			primitive.ObjectID		`json:"thingId,omitempty" bson:"thingId,omitempty"`
	Charts				[]Chart								`json:"charts" bson:"-"`
}