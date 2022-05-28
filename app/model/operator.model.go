package model

import (
	"github.com/google/uuid"
)

const TableNameOperator = "operator"

type Operator struct {
	Base
	Name           string       `gorm:"column:name;not null" json:"name"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;not null" json:"organizationId"`
	Organization   Organization `gorm:"foreignKey:OrganizationId;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ThingIds       []uuid.UUID  `gorm:"-" json:"operatorIds,omitempty"`
}

func (*Operator) TableName() string {
	return TableNameOperator
}
