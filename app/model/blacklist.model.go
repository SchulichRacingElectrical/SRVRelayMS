package model

const TableNameBlacklist = "blacklist"

type Blacklist struct {
	Base
	token      string `gorm:"column;not null;unique"`
	expiration int64  `gorm:"expiration;not null"`
}

func (*Blacklist) TableName() string {
	return TableNameChart
}
