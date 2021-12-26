package models

import "gopkg.in/mgo.v2/bson"

type Sensor struct {
	ID                   bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	SID                  int           `json:"sid,omitempty" bson:"sid,omitempty"`
	Type                 string        `json:"type,omitempty" bson:"type,omitempty"`
	LastUpdate           int64         `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Category             string        `json:"category,omitempty" bson:"category,omitempty"`
	Name                 string        `json:"name,omitempty" bson:"name,omitempty" firestore:"name,omitempty"`
	Frequency            int           `json:"frequency,omitempty" bson:"frequency,omitempty"`
	Unit                 string        `json:"unit,omitempty" bson:"unit,omitempty"`
	CanId                int           `json:"canId,omitempty" bson:"canId,omitempty"`
	Disabled             bool          `json:"disabled,omitempty" bson:"disabled,omitempty"`
	ThingID              bson.ObjectId `json:"thingId,omitempty" bson:"thingId,omitempty"`
	UpperCalibration     float64       `json:"upperCalibration,omitempty" bson:"upperCalibration,omitempty"`
	LowerCalibration     float64       `json:"lowerCalibration,omitempty" bson:"lowerCalibration,omitempty"`
	ConversionMultiplier float64       `json:"conversionMultiplier,omitempty" bson:"conversionMultiplier,omitempty"`
	UpperWarning         float64       `json:"upperWarning,omitempty" bson:"upperWarning,omitempty"`
	LowerWarning         float64       `json:"lowerWarning,omitempty" bson:"lowerWarning,omitempty"`
	UpperDanger          float64       `json:"upperDanger,omitempty" bson:"upperDanger,omitempty"`
}
