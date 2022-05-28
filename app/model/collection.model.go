package model

import (
	"github.com/google/uuid"
)

const TableNameCollection = "collection"

type Collection struct {
	Base
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description,omitempty"`
	ThingId     uuid.UUID `gorm:"type:uuid;column:thing_id;not null" json:"thingId"`
	Thing       Thing     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*Collection) TableName() string {
	return TableNameCollection
}
