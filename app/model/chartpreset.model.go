package model

import (
	"github.com/google/uuid"
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

// func (c *ChartPreset) AfterCreate() error {
// 	// Insert charts and their relationship with sensors
// 	return nil
// }

// func (c *ChartPreset) AfterFind() error {
// 	return nil
// }
