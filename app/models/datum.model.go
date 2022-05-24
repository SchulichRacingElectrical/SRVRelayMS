package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Datum struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SessionID primitive.ObjectID `json:"sessionId,omitempty" bson:"sessionId,omitempty"`
	SensorID  primitive.ObjectID `json:"sensorId,omitempty" bson:"sensorId,omitempty"`
	Value     float64            `json:"value,omitempty" bson:"value,omitempty"`
	Timestamp int64              `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}
