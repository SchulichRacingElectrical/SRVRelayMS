package model

import (
	"github.com/google/uuid"
)

const TableNameCollectionComment = "collection_comment"

type CollectionComment struct {
	Base
	CollectionId uuid.UUID  `gorm:"type:uuid;column:collection_id;not null" json:"collectionId"`
	UserId       uuid.UUID  `gorm:"type:uuid;column:user_id;not null" json:"userId"`
	LastUpdate   int64      `gorm:"column:last_update;not null" json:"lastUpdate"`
	Content      string     `gorm:"column:content;not null" json:"content"`
	Collection   Collection `gorm:"foreignKey:CollectionId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User         User       `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*CollectionComment) TableName() string {
	return TableNameCollectionComment
}
