package model

import "github.com/google/uuid"

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
