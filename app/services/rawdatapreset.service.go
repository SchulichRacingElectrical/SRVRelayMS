package services

import (
	"database-ms/config"

	"gopkg.in/mgo.v2"
)

type RawDataPresetServiceInterface interface {

}

type RawDataPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewRawDataPresetService(db *mgo.Session, c *config.Configuration) RawDataPresetServiceInterface {
	return &RawDataPresetService{db: db, config: c}
}