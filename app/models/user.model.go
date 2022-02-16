package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	DisplayName    string        `json:"name,omitempty" bson:"name,omitempty"`
	Email          string        `json:"email,omitempty" bson:"email,omitempty"`
	Password       string        `json:"password,omitempty" bson:"password,omitempty"`
	OrganizationId string        `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
	Roles          string        `json:"roles,omitempty" bson:"roles,omitempty"`
}
