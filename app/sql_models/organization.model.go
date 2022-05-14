package sqlmodels

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID        uint `gorm:"primaryKey" json:"_id,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	ApiKey    string         `json:"apiKey,omitempty"`
}
