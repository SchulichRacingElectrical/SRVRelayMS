package model

import (
	"github.com/google/uuid"
)

const TableNameThingOperator = "thing_operator"

type ThingOperator struct {
	Base
	OperatorId uuid.UUID `gorm:"type:uuid;column:operator_id;not null" json:"operatorId"`
	ThingId    uuid.UUID `gorm:"type:uuid;column:thing_id;not null" json:"thingId"`
	Operator   Operator  `gorm:"foreignKey:OperatorId;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Thing      Thing     `gorm:"foreignKey:ThingId;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*ThingOperator) TableName() string {
	return TableNameThingOperator
}
