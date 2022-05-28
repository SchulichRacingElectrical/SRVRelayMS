package model

import (
	"github.com/google/uuid"
)

const TableNameThing = "thing"

type Thing struct {
	Base
	Name           string       `gorm:"column:name;not null" json:"name"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;not null" json:"organizationId"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OperatorIds    []uuid.UUID  `gorm:"-" json:"operatorIds,omitempty"`
}

func (*Thing) TableName() string {
	return TableNameThing
}
