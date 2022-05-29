package model

import "github.com/google/uuid"

const TableNameDatum = "datum"

type Datum struct {
	Base
	Timestamp int64     `gorm:"column:timestamp;not null"`
	Value     float64   `gorm:"column:value;not null"`
	SensorId  uuid.UUID `gorm:"column:sensor_id;not null"`
	SessionId uuid.UUID `gorm:"column:session_id;not null"`
	Sensor    Sensor    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Session   Session   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (*Datum) TableName() string {
	return TableNameDatum
}
