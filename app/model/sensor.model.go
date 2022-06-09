package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameSensor = "sensor"

type LastUpdateSensors struct {
	Sensors   []Sensor    `json:"sensors"`
	SensorIds []uuid.UUID `json:"existingSensorIds"`
}

type Sensor struct {
	Base
	SmallId              int       `gorm:"column:small_id;not null;uniqueIndex:unique_sensor_smallid_in_thing" json:"smallId"`
	Type                 string    `gorm:"type:varchar(1);column:type;not null" json:"type"`
	LastUpdate           int64     `gorm:"column:last_update;not null" json:"lastUpdate"`
	Name                 string    `gorm:"column:name;not null;uniqueIndex:unique_sensor_name_in_thing" json:"name"`
	Frequency            int32     `gorm:"column:frequency;not null" json:"frequency"`
	Unit                 string    `gorm:"column:unit" json:"unit,omitempty"`
	CanId                int64     `gorm:"column:can_id;not null" json:"canId"`
	CanOffset            int       `gorm:"column:can_offset;not null" json:"canOffset"`
	ThingId              uuid.UUID `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_sensor_name_in_thing;uniqueIndex:unique_sensor_smallid_in_thing" json:"thingId"`
	UpperCalibration     float64   `gorm:"type:float;column:upper_calibration" json:"upperCalibration,omitempty"`
	LowerCalibration     float64   `gorm:"column:lower_calibration" json:"lowerCalibration,omitempty"`
	ConversionMultiplier float64   `gorm:"column:conversion_multiplier" json:"conversionMultiplier,omitempty"`
	UpperWarning         float64   `gorm:"column:upper_warning" json:"upperWarning,omitempty"`
	LowerWarning         float64   `gorm:"column:lower_warning" json:"lowerWarning,omitempty"`
	UpperDanger          float64   `gorm:"column:upper_danger" json:"upperDanger,omitempty"`
	LowerDanger          float64   `gorm:"column:lower_danger" json:"lowerDanger,omitempty"`
	UpperBound           float64   `gorm:"column:upper_bound;not null" json:"upperBound"`
	LowerBound           float64   `gorm:"column:lower_bound;not null" json:"lowerBound"`
	Thing                Thing     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*Sensor) TableName() string {
	return TableNameSensor
}

func (s *Sensor) BeforeDelete(db *gorm.DB) (err error) {
	err = db.Transaction(func(db *gorm.DB) error {
		// Delete the raw data preset entries
		var rawPresetSensors []*RawDataPresetSensor
		result := db.Table(TableNameRawdataPresetSensor).Where("sensor_id = ?", s.Id).Find(&rawPresetSensors)
		if result.Error != nil {
			return result.Error
		}
		for _, preset := range rawPresetSensors {
			err := preset.BeforeDelete(db)
			if err != nil {
				return err
			}
		}

		// Delete the chart sensor entries
		var chartSensors []*ChartSensor
		result = db.Table(TableNameChartSensor).Where("sensor_id = ?", s.Id).Find(&chartSensors)
		if result.Error != nil {
			return result.Error
		}
		for _, preset := range chartSensors {
			err := preset.BeforeDelete(db)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
