package services

import (
	"database-ms/config"

	"gopkg.in/mgo.v2"
)

type ChartPresetServiceInterface interface {

}

type ChartPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewChartPresetService(db *mgo.Session, c *config.Configuration) ChartPresetServiceInterface {
	return &ChartPresetService{db: db, config: c}
}