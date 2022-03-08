package handlers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createEntityRes struct {
	ID primitive.ObjectID `json:"_id"`
}
