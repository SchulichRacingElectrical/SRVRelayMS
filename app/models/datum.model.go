package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Datum struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SessionID primitive.ObjectID `json:"sessionId,omitempty" bson:"sessionId,omitempty"`
	SensorID  primitive.ObjectID `json:"sensorId,omitempty" bson:"sensorId,omitempty"`
	Value     float64            `json:"value" bson:"value"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
}

type FormattedDatum struct {
	X float64 // Value `json:"x`
	Y int64   // Timestamp `json:"y"`
}
