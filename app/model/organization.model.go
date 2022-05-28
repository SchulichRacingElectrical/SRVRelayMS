package model

const TableNameOrganization = "organization"

type Organization struct {
	Base
	Name   string `gorm:"column:name;not null;unique" json:"name"`
	APIKey string `gorm:"column:api_key;not null;unique" json:"apiKey,omitempty"`
}

func (*Organization) TableName() string {
	return TableNameOrganization
}
