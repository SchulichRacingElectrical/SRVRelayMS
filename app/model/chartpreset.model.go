package model

import (
	"github.com/google/uuid"
)

const TableNameChartpreset = "chartpreset"

type ChartPreset struct {
	Base
	Name    string    `gorm:"column:name;not null" json:"name"`
	ThingId uuid.UUID `gorm:"type:uuid;column:thing_id;not null" json:"thingId"`
	Thing   Thing     `gorm:"foreignKey:ThingId;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*ChartPreset) TableName() string {
	return TableNameChartpreset
}
