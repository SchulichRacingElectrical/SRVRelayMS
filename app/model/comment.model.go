package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameComment = "comment"

type Comment struct {
	Base
	CollectionId *uuid.UUID `gorm:"type:uuid;column:collection_id" json:"collectionId,omitempty"`
	SessionId    *uuid.UUID `gorm:"type:uuid;column:session_id" json:"sessionId,omitempty"`
	UserId       uuid.UUID  `gorm:"type:uuid;column:user_id" json:"userId"`
	CommentId    *uuid.UUID `gorm:"type:uuid;column:comment_id" json:"commentId,omitempty"`
	Username     string     `gorm:"column:username" json:"username"`
	LastUpdate   int64      `gorm:"column:last_update;not null" json:"lastUpdate"`
	Content      string     `gorm:"column:content; not null" json:"content"`
	Collection   Collection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Session      Session    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User         User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Comments     []Comment  `gorm:"foreignKey:comment_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
}

func (*Comment) TableName() string {
	return TableNameComment
}

func (c *Comment) AfterFind(db *gorm.DB) (err error) {
	if c.CommentId == nil {
		c.Comments = []Comment{}
	} else {
		var replies []Comment
		result := db.Order("last_update desc").Find(&replies, "comment_id = ?", c.Id)
		if result.Error != nil {
			return result.Error
		}
		c.Comments = replies
	}
	return nil
}
