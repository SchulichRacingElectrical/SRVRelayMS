package model

import (
	"github.com/google/uuid"
)

const TableNameComment = "comment"

type Comment struct {
	Base
	CollectionId *uuid.UUID `gorm:"type:uuid;column:collection_id" json:"collectionId,omitempty"`
	SessionId    *uuid.UUID `gorm:"type:uuid;column:session_id" json:"sessionId,omitempty"`
	UserId       uuid.UUID  `gorm:"type:uuid;column:user_id" json:"userId"`
	LastUpdate   int64      `gorm:"column:last_update;not null" json:"lastUpdate"`
	Content      string     `gorm:"column:content; not null" json:"content"`
	Collection   Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Session      Session    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User         User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*Comment) TableName() string {
	return TableNameComment
}
