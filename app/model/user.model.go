package model

import (
	"github.com/google/uuid"
)

const TableNameUser = "user"

type User struct {
	Base
	DisplayName    string       `gorm:"column:display_name;not null;uniqueIndex:unique_name_in_org" json:"name"`
	Email          string       `gorm:"column:email;not null;unique" json:"email"`
	Password       string       `gorm:"column:password;not null" json:"-"`
	OrganizationId uuid.UUID    `gorm:"type:uuid;column:organization_id;uniqueIndex:unique_name_in_org" json:"organizationId"`
	Role           string       `gorm:"column:role;not null" json:"role"`
	Organization   Organization `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (*User) TableName() string {
	return TableNameUser
}
