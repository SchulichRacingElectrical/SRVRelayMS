package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameRawdataPresetSensor = "rawdatapreset_sensor"

type RawDataPresetSensor struct {
	Base
	RawDataPresetId uuid.UUID     `gorm:"type:uuid;column:rawdatapreset_id;not null"`
	SensorId        uuid.UUID     `grom:"type:uuid;column:sensor_id;not null"`
	RawDataPreset   RawDataPreset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Sensor          Sensor        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*RawDataPresetSensor) TableName() string {
	return TableNameRawdataPresetSensor
}

func (rs *RawDataPresetSensor) BeforeDelete(db *gorm.DB) (err error) {
	// Find all the other sensors attached to the chart
	var presetSensors []*RawDataPresetSensor
	result := db.Table(TableNameRawdataPresetSensor).Where("rawdatapreset_id = ?", rs.RawDataPresetId).Find(&presetSensors)
	if result.Error != nil {
		return result.Error
	}

	// Delete the preset if this is the last sensor attached to it
	if len(presetSensors) == 1 {
		preset := RawDataPreset{Base: Base{Id: rs.RawDataPresetId}}
		result := db.Table(TableNameRawdatapreset).Delete(&preset)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
