package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RunFileUpload struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ThingID         primitive.ObjectID `json:"thingId,omitempty" bson:"thingId,omitempty" form:"thingId"`
	RunId           primitive.ObjectID `json:"runId,omitempty" bson:"runId,omitempty" form:"runId"`
	OperatorId      primitive.ObjectID `json:"operatorId,omitempty" bson:"operatorId,omitempty" form:"operatorId"`
	UploadDateEpoch int64              `json:"uploadDateEpoch,omitempty" bson:"uploadDateEpoch,omitempty"`
}

type RunFileMetaData struct {
	FileName string `json:"filename,omitempty" bson:"filename,omitempty"`
}
