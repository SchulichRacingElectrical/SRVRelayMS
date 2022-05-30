package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (cs *ChartSensor) BeforeDelete(db *gorm.DB) (err error) {
	// Find all the other sensors attached to the chart
	var chartSensors []*ChartSensor
	result := db.Table(TableNameChartSensor).Where("chart_id = ?", cs.ChartId).Find(&chartSensors)
	if result.Error != nil {
		return result.Error
	}

	// Delete the chart if this is the last sensor attached to it
	if len(chartSensors) == 1 {
		chart := Chart{Base: Base{Id: cs.ChartId}}
		result := db.Delete(&chart)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
