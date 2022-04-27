package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thing struct {
	ID        				primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      				string               `json:"name,omitempty" bson:"name,omitempty"`
	OrganizationId		primitive.ObjectID	 `json:"organizationId,omitempty" bson:"organizationId,omitempty"` 	
	OperatorIds				[]primitive.ObjectID `json:"operatorIds,omitempty" bson:"-"`	
	// TODO: Add JSON field
}
