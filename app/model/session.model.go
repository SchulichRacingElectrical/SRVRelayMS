package model

import (
	"github.com/google/uuid"
)

const TableNameSession = "session"

type Session struct {
	Base
	Name         string     `gorm:"column:name;not null" json:"name"`
	StartTime    int64      `gorm:"column:start_time;not null" json:"startTime"`
	EndTime      int64      `gorm:"column:end_time" json:"endTime,omitempty"`
	CollectionId uuid.UUID  `gorm:"type:uuid;column:collection_id" json:"collectionId,omitempty"`
	ThingId      uuid.UUID  `gorm:"type:uuid;column:thing_id;not null" json:"thingId"`
	Collection   Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Thing        Thing      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*Session) TableName() string {
	return TableNameSession
}
