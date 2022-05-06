package services

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"

	"gopkg.in/mgo.v2"
)

type ChartPresetServiceInterface interface {
	Create(context.Context, *models.ChartPreset) error
	FindByThingId(context.Context, string) ([]*models.ChartPreset, error)
	Update(context.Context, *models.ChartPreset) error
	Delete(context.Context, string) error
	FindById(context.Context, string) (*models.ChartPreset, error)
	AreChartsValid(context.Context, *models.ChartPreset) bool
	IsChartPresetUnique(context.Context, *models.ChartPreset) bool
}

type ChartPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewChartPresetService(db *mgo.Session, c *config.Configuration) ChartPresetServiceInterface {
	return &ChartPresetService{db: db, config: c}
}

func (service *ChartPresetService) Create(ctx context.Context, chartPreset *models.ChartPreset) error {
	return nil
}

func (service *ChartPresetService) FindByThingId(ctx context.Context, thingId string) ([]*models.ChartPreset, error) {
	return nil, nil
}

func (service *ChartPresetService) Update(ctx context.Context, updatedChartPreset *models.ChartPreset) error {
	return nil
}

func (service *ChartPresetService) Delete(ctx context.Context, chartPresetId string) error {
	return nil
}

func (service *ChartPresetService) FindById(ctx context.Context, chartPresetId string) (*models.ChartPreset, error) {
	return nil, nil
}

func (service *ChartPresetService) AreChartsValid(ctx context.Context, chartPreset *models.ChartPreset) bool {
	return false
}

func (service *ChartPresetService) IsChartPresetUnique(ctx context.Context, chartPreset *models.ChartPreset) bool {
	return false
}



