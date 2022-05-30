package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameSensor = "sensor"

type Sensor struct {
	Base
	SmallId              int       `gorm:"column:small_id;not null;uniqueIndex:unique_sensor_smallid_in_thing" json:"smallId"`
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

func (s *Sensor) BeforeDelete(db *gorm.DB) error {
	// Find the raw data ids associated with the sensor
	var presetSensors []*RawDataPresetSensor
	result := db.Table(TableNameRawdataPresetSensor).Where("sensor_id = ?", s.Id).Find(&presetSensors)
	if result.Error != nil {
		return result.Error
	}

	// For each raw data preset, delete it if there are no sensors remaining
	for _, presetSensor := range presetSensors {
		var allPresetSensors []*RawDataPresetSensor
		result := db.Table(TableNameRawdataPresetSensor).Where("rawdatapreset_id", presetSensor.RawDataPresetId).Find(&allPresetSensors)
		if result.Error != nil {
			return result.Error
		}

		// Delete the preset if the sensor about to be deleted is the last one in the preset
		if len(allPresetSensors) == 1 {
			if allPresetSensors[1].SensorId == s.Id {
				preset := &RawDataPreset{Base: Base{Id: presetSensor.RawDataPresetId}}
				result := db.Delete(&preset)
				if result.Error != nil {
					return result.Error
				}
			}
		}
	}

	// Find the chart ids associated with the sensor
	var chartSensors []*ChartSensor
	result = db.Table(TableNameChartSensor).Where("sensor_id = ?", s.Id).Find(&chartSensors)
	if result.Error != nil {
		return result.Error
	}

	// For each chart, delete it if there are no sensors remaining
	var presetsToCheck []uuid.UUID
	for _, chartSensor := range chartSensors {
		var allChartSensors []*ChartSensor
		result := db.Table(TableNameChartSensor).Where("chart_id", chartSensor.ChartId).Find(&allChartSensors)
		if result.Error != nil {
			return result.Error
		}

		// Delete the chart if the sensor about to be deleted is the last one in the chart
		if len(allChartSensors) == 1 {
			if allChartSensors[1].SensorId == s.Id {
				chart := &Chart{Base: Base{Id: chartSensor.ChartId}}

				// Find the chart and add its preset to the list to check
				result := db.First(&chart)
				if result.Error != nil {
					return result.Error
				}
				presetsToCheck = append(presetsToCheck, chart.ChartPresetId)

				// Delete the chart
				result = db.Delete(&chart)
				if result.Error != nil {
					return result.Error
				}
			}
		}
	}

	// Check all of the chartPresets and check if they need to be deleted (if there are no charts left)
	for _, chartPresetId := range presetsToCheck {
		// Find the preset
		preset := ChartPreset{Base: Base{Id: chartPresetId}}
		result := db.Find(&preset)
		if result.Error != nil {
			continue
		}

	}
	return nil
}
