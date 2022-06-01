package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (t *Thing) AfterCreate(db *gorm.DB) (err error) {
	return InsertThingOperators(t, db)
}

func (t *Thing) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all of the associated thing-operators
	result := db.Table(TableNameThingOperator).Where("thing_id = ?", t.Id).Delete(&ThingOperator{})
	if result.Error != nil {
		return result.Error
	}

	// Write the new thing-operators
	return InsertThingOperators(t, db)
}

func (t *Thing) AfterFind(db *gorm.DB) (err error) {
	var thingOperators []*ThingOperator
	t.OperatorIds = []uuid.UUID{}
	result := db.Table(TableNameThingOperator).Where("thing_id = ?", t.Id).Find(&thingOperators)
	if result.Error != nil {
		return result.Error
	}
	for _, thingOperator := range thingOperators {
		t.OperatorIds = append(t.OperatorIds, thingOperator.OperatorId)
	}
	return nil
}

func InsertThingOperators(t *Thing, db *gorm.DB) (err error) {
	// Generate the list of thing-operators
	thingOperators := []ThingOperator{}
	for _, operatorId := range t.OperatorIds {
		thingOperator := ThingOperator{}
		thingOperator.ThingId = t.Id
		thingOperator.OperatorId = operatorId
		thingOperators = append(thingOperators, thingOperator)
	}

	// Insert empty operatorIds
	if len(t.OperatorIds) == 0 {
		t.OperatorIds = []uuid.UUID{}
		return
	}

	// Batch insert thing-operators
	result := db.Table(TableNameThingOperator).CreateInBatches(thingOperators, 100)
	return result.Error
}
