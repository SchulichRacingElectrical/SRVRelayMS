package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChartPreset struct {
	ID					primitive.ObjectID		`json:"_id,omitempty" bson:"_id,omitempty"`
}