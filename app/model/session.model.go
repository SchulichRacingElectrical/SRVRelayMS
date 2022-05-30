package model

import (
	"github.com/google/uuid"
)

const TableNameSession = "session"

type Session struct {
	Base
	Name         string     `gorm:"column:name;not null;uniqueIndex:unique_session_name_in_thing" json:"name"`
	StartTime    int64      `gorm:"column:start_time;not null" json:"startTime"`
	EndTime      int64      `gorm:"column:end_time" json:"endTime,omitempty"`
	FileName     string     `gorm:"column:file_name;unique" json:"fileName,omitempty"`
	CollectionId uuid.UUID  `gorm:"type:uuid;column:collection_id" json:"collectionId,omitempty"`
	ThingId      uuid.UUID  `gorm:"type:uuid;column:thing_id;not null;uniqueIndex:unique_session_name_in_thing" json:"thingId"`
	Collection   Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Thing        Thing      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*Session) TableName() string {
	return TableNameSession
}
