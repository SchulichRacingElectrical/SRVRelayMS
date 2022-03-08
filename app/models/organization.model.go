package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	ApiKey string             `json:"api_key,omitempty" bson:"api_key,omitempty"`
}
