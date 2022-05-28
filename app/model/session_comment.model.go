package model

import (
	"github.com/google/uuid"
)

const TableNameSessionComment = "session_comment"

type SessionComment struct {
	Base
	SessionId  uuid.UUID `gorm:"type:uuid;column:session_id;not null" json:"sessionId"`
	UserId     uuid.UUID `gorm:"type:uuid;column:user_id;not null" json:"userId"`
	LastUpdate int64     `gorm:"column:last_update;not null" json:"lastUpdate"`
	Content    string    `gorm:"column:content;not null" json:"content"`
	Session    Session   `gorm:"foreignKey:SessionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User       User      `gorm:"foreignKey:UserId;references:Id;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*SessionComment) TableName() string {
	return TableNameSessionComment
}
