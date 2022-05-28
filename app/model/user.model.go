package model

import (
	"github.com/google/uuid"
)

const TableNameUser = "user"

type User struct {
	Base
	DisplayName    string       `gorm:"column:display_name;not null" json:"displayName"`
	Email          string       `gorm:"column:email;not null" json:"email"`
	Password       string       `gorm:"column:password;not null" json:"password,omitempty"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id" json:"organizationId"`
	Role           string       `gorm:"column:role;not null" json:"role"`
	Organization   Organization `gorm:"foreignKey:OrganizationId;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (*User) TableName() string {
	return TableNameUser
}
