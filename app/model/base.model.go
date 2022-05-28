package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id uuid.UUID `gorm:"type:uuid;column:id;primaryKey;" json:"_id"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewString()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
