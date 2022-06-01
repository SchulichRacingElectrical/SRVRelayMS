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
	Session    Session   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*SessionComment) TableName() string {
	return TableNameSessionComment
}
