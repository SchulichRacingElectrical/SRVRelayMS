package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Organization struct {
	ID     bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name,omitempty" bson:"name,omitempty"`
	ApiKey string        `json:"api_key,omitempty" bson:"api_key,omitempty"`
}
