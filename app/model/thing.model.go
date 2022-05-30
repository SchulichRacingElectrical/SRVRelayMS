package model

import (
	"github.com/google/uuid"
)

const TableNameThing = "thing"

type Thing struct {
	Base
	Name           string       `gorm:"column:name;not null;uniqueIndex:unique_thing_name_in_org" json:"name"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;not null;uniqueIndex:unique_thing_name_in_org" json:"organizationId"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	OperatorIds    []uuid.UUID  `gorm:"-" json:"operatorIds"`
}

func (*Thing) TableName() string {
	return TableNameThing
}
