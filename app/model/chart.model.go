package model

import (
	"github.com/google/uuid"
)

const TableNameChart = "chart"

type Chart struct {
	Base
	Name          string      `gorm:"column:name;not null" json:"name"`
	Type          string      `gorm:"column:type;not null" json:"type"`
	ChartPresetId uuid.UUID   `gorm:"type:uuid;column:chart_preset_id" json:"chartPresetId"`
	SensorIds     []uuid.UUID `gorm:"type:uuid[];column:sensor_ids;not null" json:"sensorIds"`
	ChartPreset   ChartPreset `gorm:"foreignKey:ChartPresetId;references:Id;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*Chart) TableName() string {
	return TableNameChart
}
