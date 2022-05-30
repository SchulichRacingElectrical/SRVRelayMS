package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameOperator = "operator"

type Operator struct {
	Base
	Name           string       `gorm:"column:name;not null;uniqueIndex:unique_operator_name_in_org" json:"name"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;not null;uniqueIndex:unique_operator_name_in_org" json:"organizationId"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	ThingIds       []uuid.UUID  `gorm:"-" json:"thingIds"`
	// Should we just send all the things back?
}

func (*Operator) TableName() string {
	return TableNameOperator
}

func (o *Operator) AfterCreate(db *gorm.DB) (err error) {
	return InsertOperatorThings(o, db)
}

func (o *Operator) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all of the associated thing operators
	result := db.Table(TableNameThingOperator).Where("operator_id = ?", o.Id).Delete(&ThingOperator{})
	if result.Error != nil {
		return result.Error
	}

	// Write the new operator-things
	return InsertOperatorThings(o, db)
}

func (o *Operator) AfterFind(db *gorm.DB) (err error) {
	// Get the ids of the relationship with thing
	var thingOperators []*ThingOperator
	o.ThingIds = []uuid.UUID{}
	result := db.Table(TableNameThingOperator).Where("operator_id = ?", o.Id).Find(&thingOperators)
	if result.Error != nil {
		return result.Error
	}
	for _, thingOperator := range thingOperators {
		o.ThingIds = append(o.ThingIds, thingOperator.ThingId)
	}
	return nil
}

func InsertOperatorThings(o *Operator, db *gorm.DB) (err error) {
	// Regenerate the list of thing-operators
	var thingOperators []ThingOperator
	for _, thingId := range o.ThingIds {
		thingOperators = append(thingOperators, ThingOperator{
			OperatorId: o.Id,
			ThingId:    thingId,
		})
	}

	// Insert empty thingIds
	if len(o.ThingIds) == 0 {
		o.ThingIds = []uuid.UUID{}
	}

	// Batch insert thing-operators
	result := db.Table(TableNameThingOperator).CreateInBatches(thingOperators, 100)
	return result.Error
}
