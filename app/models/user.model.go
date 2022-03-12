package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DisplayName    string             `json:"name,omitempty" bson:"name,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	Password       string             `json:"password,omitempty" bson:"password,omitempty"`
	OrganizationId primitive.ObjectID `json:"organizationId,omitempty" bson:"organizationId,omitempty"`
	// These roles are: Admin, Lead, Member, Guest, Pending
	Roles string `json:"roles,omitempty" bson:"roles,omitempty"`
}
