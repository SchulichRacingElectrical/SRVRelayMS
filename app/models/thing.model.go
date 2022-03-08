package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thing struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty" firestore:"name,omitempty"`
	Operators []primitive.ObjectID `json:"operators,omitempty" bson:"operators,omitempty"`
	Sensors   []primitive.ObjectID `json:"sensors,omitempty" bson:"sensors,omitempty"`
}

type ThingUpdate struct {
	Name string `json:"name,omitempty" bson:"name,omitempty" firestore:"name,omitempty"`
}
