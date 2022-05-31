package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameChart = "chart"

type Chart struct {
	Base
	Name          string      `gorm:"column:name;not null;uniqueIndex:unique_chart_name_in_preset" json:"name"`
	Type          string      `gorm:"column:type;not null" json:"type"`
	ChartPresetId uuid.UUID   `gorm:"type:uuid;column:chart_preset_id;uniqueIndex:unique_chart_name_in_preset" json:"chartPresetId,omitempty"`
	ChartPreset   ChartPreset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SensorIds     []uuid.UUID `gorm:"-" json:"sensorIds"`
}

func (*Chart) TableName() string {
	return TableNameChart
}

func (c *Chart) AfterCreate(db *gorm.DB) (err error) {
	return InsertChartSensors(c, db)
}

func (c *Chart) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all the associated chart-sensors
	result := db.Table(TableNameChartSensor).Where("chart_id = ?", c.Id).Delete(&ChartSensor{})
	if result.Error != nil {
		return result.Error
	}

	// Write the new chart-sensors
	return InsertChartSensors(c, db)
}

func (c *Chart) AfterFind(db *gorm.DB) (err error) {
	var chartSensors []*ChartSensor
	c.SensorIds = []uuid.UUID{}
	result := db.Table(TableNameChartSensor).Where("chart_id = ?", c.Id).Find(&chartSensors)
	if result.Error != nil {
		return result.Error
	}
	for _, chartSensor := range chartSensors {
		c.SensorIds = append(c.SensorIds, chartSensor.SensorId)
	}
	return nil
}

func (c *Chart) BeforeDelete(db *gorm.DB) (err error) {
	// Find all of the associated charts to the preset
	var allCharts []*Chart
	result := db.Table(TableNameChart).Where("chart_preset_id = ?", c.ChartPresetId).Find(&allCharts)
	if result.Error != nil {
		return result.Error
	}

	// Delete the preset if the only chart remaining is c
	if len(allCharts) == 1 {
		chartPreset := ChartPreset{Base: Base{Id: c.ChartPresetId}}
		result := db.Delete(&chartPreset)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func InsertChartSensors(c *Chart, db *gorm.DB) (err error) {
	// Generate the new list of chart-sensors
	chartSensors := []ChartSensor{}
	for _, sensorId := range c.SensorIds {
		chartSensor := ChartSensor{}
		chartSensor.ChartId = c.Id
		chartSensor.SensorId = sensorId
		chartSensors = append(chartSensors, chartSensor)
	}

	// Insert empty sensorIds
	if len(c.SensorIds) == 0 {
		c.SensorIds = []uuid.UUID{}
		return
	}

	// Batch insert the chart-sensors
	result := db.Table(TableNameChartSensor).CreateInBatches(chartSensors, 100)
	return result.Error
}
