package model

import "github.com/google/uuid"

const TableNameChartSensor = "chart_sensor"

type ChartSensor struct {
	Base
	ChartId  uuid.UUID `gorm:"type:uuid;column:chart_id;not null"`
	SensorId uuid.UUID `gorm:"type:uuid;column:sensor_id;not null"`
	Chart    Chart     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Sensor   Sensor    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*ChartSensor) TableName() string {
	return TableNameChartSensor
}
