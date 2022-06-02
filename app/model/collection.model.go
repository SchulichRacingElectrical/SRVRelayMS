package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameCollection = "collection"

type Collection struct {
	Base
	Name        string      `gorm:"column:name;not null;uniqueIndex:unique_collection_name_in_thing" json:"name"`
	Description string      `gorm:"column:description" json:"description,omitempty"`
	ThingId     uuid.UUID   `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_collection_name_in_thing" json:"thingId"`
	Thing       Thing       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SessionIds  []uuid.UUID `gorm:"-" json:"sessionIds"`
}

func (*Collection) TableName() string {
	return TableNameCollection
}

func (c *Collection) AfterCreate(db *gorm.DB) (err error) {
	return InsertCollectionSessions(c, db)
}

func (c *Collection) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all of the associated collection-sessions
	result := db.Table(TableNameSessionCollection).Where("collection_id = ?", c.Id).Delete(&SessionCollection{})
	if result.Error != nil {
		return result.Error
	}

	// Write the new collection-sessions
	return InsertCollectionSessions(c, db)
}

func (c *Collection) AfterFind(db *gorm.DB) (err error) {
	var collectionSessions []*SessionCollection
	c.SessionIds = []uuid.UUID{}
	result := db.Table(TableNameSessionCollection).Where("collection_id = ?", c.Id).Find(&collectionSessions)
	if result.Error != nil {
		return result.Error
	}
	for _, collectionSession := range collectionSessions {
		c.SessionIds = append(c.SessionIds, collectionSession.SessionId)
	}
	return nil
}

func InsertCollectionSessions(c *Collection, db *gorm.DB) (err error) {
	// Insert empty sessionIds
	if len(c.SessionIds) == 0 {
		c.SessionIds = []uuid.UUID{}
		return
	}

	// Generate the list of collection-sessions
	collectionSessions := []SessionCollection{}
	for _, sessionId := range c.SessionIds {
		collectionSession := SessionCollection{}
		collectionSession.CollectionId = c.Id
		collectionSession.SessionId = sessionId
		collectionSessions = append(collectionSessions, collectionSession)
	}

	// Batch insert thing-operators
	result := db.Table(TableNameSessionCollection).CreateInBatches(collectionSessions, 100)
	return result.Error
}
