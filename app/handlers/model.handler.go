package handlers

import (
	"gopkg.in/mgo.v2/bson"
)

type createSensorRes struct {
	ID bson.ObjectId `json:"_id"`
}

type createOrganizationRes struct {
	ID bson.ObjectId `json:"_id"`
}

type createUserRes struct {
	ID bson.ObjectId `json:"_id"`
}
