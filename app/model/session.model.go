package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameSession = "session"

type Session struct {
	Base
	Name          string      `gorm:"column:name;not null;uniqueIndex:unique_session_name_in_thing" json:"name"`
	StartTime     int64       `gorm:"column:start_time;not null" json:"startTime"`
	EndTime       int64       `gorm:"column:end_time" json:"endTime,omitempty"`
	FileName      string      `gorm:"column:file_name;unique" json:"fileName,omitempty"`
	ThingId       uuid.UUID   `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_session_name_in_thing" json:"thingId"`
	OperatorId    *uuid.UUID  `gorm:"type:uuid;column:operator_id" json:"operatorId,omitempty"`
	Thing         Thing       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Operator      Operator    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	CollectionIds []uuid.UUID `gorm:"-" json:"collectionIds"`
}

func (*Session) TableName() string {
	return TableNameSession
}

func (s *Session) AfterCreate(db *gorm.DB) (err error) {
	return InsertSessionCollections(s, db)
}

func (s *Session) AfterUpdate(db *gorm.DB) (err error) {
	// Delete all of the associated session-collections
	result := db.Table(TableNameSessionCollection).Where("session_id = ?", s.Id).Delete(&SessionCollection{})
	if result.Error != nil {
		return result.Error
	}

	// Write the new session-collections
	return InsertSessionCollections(s, db)
}

func (s *Session) AfterFind(db *gorm.DB) (err error) {
	var sessionCollections []*SessionCollection
	s.CollectionIds = []uuid.UUID{}
	result := db.Table(TableNameSessionCollection).Where("session_id = ?", s.Id).Find(&sessionCollections)
	if result.Error != nil {
		return result.Error
	}
	for _, sessionCollection := range sessionCollections {
		s.CollectionIds = append(s.CollectionIds, sessionCollection.CollectionId)
	}
	return nil
}

func InsertSessionCollections(s *Session, db *gorm.DB) (err error) {
	// Insert empty sessionIds
	if len(s.CollectionIds) == 0 {
		s.CollectionIds = []uuid.UUID{}
		return
	}

	// Generate the list of collection-sessions
	sessionCollections := []SessionCollection{}
	for _, sessionId := range s.CollectionIds {
		sessionCollection := SessionCollection{}
		sessionCollection.CollectionId = s.Id
		sessionCollection.SessionId = sessionId
		sessionCollections = append(sessionCollections, sessionCollection)
	}

	// Batch insert thing-operators
	result := db.Table(TableNameSessionCollection).CreateInBatches(sessionCollections, 100)
	return result.Error
}
