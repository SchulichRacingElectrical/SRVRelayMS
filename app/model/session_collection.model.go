package model

import "github.com/google/uuid"

const TableNameSessionCollection = "session_collection"

type SessionCollection struct {
	Base
	SessionId    uuid.UUID  `gorm:"type:uuid;column:session_id;not null"`
	CollectionId uuid.UUID  `gorm:"type:uuid;column:collection_id; not null"`
	Session      Session    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Collection   Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*SessionCollection) TableName() string {
	return TableNameSessionCollection
}
