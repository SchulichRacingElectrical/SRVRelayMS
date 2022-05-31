package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameChartpreset = "chartpreset"

type ChartPreset struct {
	Base
	Name    string    `gorm:"column:name;not null;uniqueIndex:unique_chartpreset_name_in_thing" json:"name"`
	ThingId uuid.UUID `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_chartpreset_name_in_thing" json:"thingId"`
	Thing   Thing     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Charts  []Chart   `gorm:"-" json:"charts"`
}

func (*ChartPreset) TableName() string {
	return TableNameChartpreset
}

func (c *ChartPreset) AfterCreate(db *gorm.DB) (err error) {
	for i := range c.Charts {
		c.Charts[i].ChartPresetId = c.Id
	}
	result := db.Table(TableNameChart).CreateInBatches(c.Charts, 100)
	return result.Error
}

func (c *ChartPreset) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all the associated charts
	result := db.Table(TableNameChart).Where("chart_preset_id = ?", c.Id).Delete(&Chart{})
	if result.Error != nil {
		return result.Error
	}

	// Insert the new charts
	for i := range c.Charts {
		c.Charts[i].ChartPresetId = c.Id
	}
	result = db.Table(TableNameChart).CreateInBatches(c.Charts, 100)
	return result.Error
}

func (c *ChartPreset) AfterFind(db *gorm.DB) (err error) {
	c.Charts = []Chart{}
	result := db.Table(TableNameChart).Where("chart_preset_id = ?", c.Id).Find(&c.Charts)
	return result.Error
}
