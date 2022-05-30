package model

import (
	"github.com/google/uuid"
)

const TableNameOperator = "operator"

type Operator struct {
	Base
	Name           string       `gorm:"column:name;not null;uniqueIndex:unique_name_in_org" json:"name"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;not null;uniqueIndex:unique_name_in_org" json:"organizationId"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ThingIds       []uuid.UUID  `gorm:"-" json:"thingIds,omitempty"`
}

func (*Operator) TableName() string {
	return TableNameOperator
}
