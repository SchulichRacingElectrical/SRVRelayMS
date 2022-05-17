package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RawDataPreset struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty"`
	SensorIds []primitive.ObjectID `json:"sensorIds,omitempty" bson:"sensorIds,omitempty"`
	ThingId   primitive.ObjectID   `json:"thingId,omitempty" bson:"thingId,omitempty"`
}
