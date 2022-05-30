package model

import (
	"github.com/google/uuid"
)

const TableNameSensor = "sensor"

type Sensor struct {
	Base
	SmallId              int32     `gorm:"column:small_id;not null;uniqueIndex:unique_sensor_smallid_in_thing" json:"smallId"`
	Type                 string    `gorm:"type:varchar(1);column:type;not null" json:"type"`
	LastUpdate           int64     `gorm:"column:last_update;not null" json:"lastUpdate"`
	Name                 string    `gorm:"column:name;not null;uniqueIndex:unique_sensor_name_in_thing" json:"name"`
	Frequency            int32     `gorm:"column:frequency;not null" json:"frequency"`
	Unit                 string    `gorm:"column:unit" json:"unit,omitempty"`
	CanID                int64     `gorm:"column:can_id;not null;uniqueIndex:unique_sensor_can_in_thing" json:"canId"`
	ThingId              uuid.UUID `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_sensor_name_in_thing;uniqueIndex:unique_sensor_can_in_thing;uniqueIndex:unique_sensor_smallid_in_thing" json:"thingId"`
	UpperCalibration     float64   `gorm:"column:upper_calibration" json:"upperCalibration,omitempty"`
	LowerCalibration     float64   `gorm:"column:lower_calibration" json:"lowerCalibration,omitempty"`
	ConversionMultiplier float64   `gorm:"column:conversion_multiplier" json:"conversionMultiplier,omitempty"`
	UpperWarning         float64   `gorm:"column:upper_warning" json:"upperWarning,omitempty"`
	LowerWarning         float64   `gorm:"column:lower_warning" json:"lowerWarning,omitempty"`
	UpperDanger          float64   `gorm:"column:upper_danger" json:"upperDanger,omitempty"`
	LowerDanger          float64   `gorm:"column:lower_danger" json:"lowerDanger,omitempty"`
	UpperBound           float64   `gorm:"column:upper_bound;not null" json:"upperBound"`
	LowerBound           float64   `gorm:"column:lower_bound;not null" json:"lowerBound"`
	Significance         float64   `gorm:"column:significance" json:"significance,omitempty"`
	Thing                Thing     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*Sensor) TableName() string {
	return TableNameSensor
}
