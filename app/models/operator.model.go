package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Operator struct {
	ID     primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string               `json:"name,omitempty" bson:"name,omitempty" firestore:"name,omitempty"`
	Things []primitive.ObjectID `json:"things,omitempty" bson:"things,omitempty"` // TODO Do we need this?
	// Json
}
