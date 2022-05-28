package model

import (
	"github.com/google/uuid"
)

const TableNameChart = "chart"

type Chart struct {
	Base
	Name          string      `gorm:"column:name;not null;uniqueIndex:unique_name_in_preset" json:"name"`
	Type          string      `gorm:"column:type;not null" json:"type"`
	ChartPresetId uuid.UUID   `gorm:"type:uuid;column:chart_preset_id;uniqueIndex:unique_name_in_preset" json:"chartPresetId"`
	SensorIds     []uuid.UUID `gorm:"type:uuid[];column:sensor_ids;not null" json:"sensorIds"`
	ChartPreset   ChartPreset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*Chart) TableName() string {
	return TableNameChart
}
