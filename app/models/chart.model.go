package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chart struct {
	ID						primitive.ObjectID			`json:"_id,omitempty" bson:"_id,omitempty"`
	Name					string									`json:"name,omitempty" bson:"name,omitempty"`
	ChartPresetID	primitive.ObjectID			`json:"chartPresetId,omitempty" bson:"chartPresetId,omitempty"`
	SensorIds			[]primitive.ObjectID		`json:"sensorIds,omitempty" bson:"sensorIds,omitempty"`
}