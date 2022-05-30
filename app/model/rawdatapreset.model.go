package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameRawdatapreset = "rawdatapreset"

type RawDataPreset struct {
	Base
	Name      string      `gorm:"column:name;not null;uniqueIndex:unique_rawdatapreset_name_in_thing" json:"name"`
	ThingId   uuid.UUID   `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_rawdatapreset_name_in_thing" json:"thingId"`
	Thing     Thing       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SensorIds []uuid.UUID `gorm:"-" json:"sensorIds"`
}

func (*RawDataPreset) TableName() string {
	return TableNameRawdatapreset
}

func (r *RawDataPreset) AfterCreate(db *gorm.DB) (err error) {
	// Generate the list of preset sensors
	var rawDataPresetSensors []RawDataPresetSensor
	for _, sensorId := range r.SensorIds {
		rawDataPresetSensors = append(rawDataPresetSensors, RawDataPresetSensor{
			RawDataPresetId: r.Id,
			SensorId:        sensorId,
		})
	}

	// Insert empty sensorIds
	if len(r.SensorIds) == 0 {
		r.SensorIds = []uuid.UUID{}
	}

	// Batch insert preset sensors
	result := db.Table(TableNameRawdataPresetSensor).CreateInBatches(rawDataPresetSensors, 100)
	return result.Error
}

func (r *RawDataPreset) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all of the associated preset-sensors
	result := db.Table(TableNameRawdataPresetSensor).Where("rawdatapreset_id = ?", r.Id).Delete(&RawDataPresetSensor{})
	if result.Error != nil {
		return result.Error
	}

	// Generate the list of preset-sensors
	presetSensors := []RawDataPresetSensor{}
	for _, sensorId := range r.SensorIds {
		presetSensor := RawDataPresetSensor{}
		presetSensor.RawDataPresetId = r.Id
		presetSensor.SensorId = sensorId
		presetSensors = append(presetSensors, presetSensor)
	}

	// Batch insert preset-sensors
	result = db.Table(TableNameRawdataPresetSensor).CreateInBatches(presetSensors, 100)
	return result.Error
}

func (r *RawDataPreset) AfterFind(db *gorm.DB) (err error) {
	// Insert the associated sensor Ids
	var presetSensors []*RawDataPresetSensor
	r.SensorIds = []uuid.UUID{}
	result := db.Table(TableNameRawdataPresetSensor).Where("rawdatapreset_id = ?", r.Id).Find(&presetSensors)
	if result.Error != nil {
		return result.Error
	}
	for _, presetSensor := range presetSensors {
		r.SensorIds = append(r.SensorIds, presetSensor.SensorId)
	}
	return nil
}
