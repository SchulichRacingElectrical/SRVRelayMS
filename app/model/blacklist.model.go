package model

const TableNameBlacklist = "blacklist"

type Blacklist struct {
	Base
	Token      string `gorm:"column;column:token;not null;unique"`
	Expiration int64  `gorm:"expiration;column:expiration;not null"`
}

func (*Blacklist) TableName() string {
	return TableNameBlacklist
}
