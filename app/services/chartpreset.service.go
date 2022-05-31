package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type ChartPresetServiceInterface interface {
	// Public
	FindByThingId(context.Context, uuid.UUID) ([]*model.ChartPreset, *pgconn.PgError)
	Create(context.Context, *model.ChartPreset) *pgconn.PgError
	Update(context.Context, *model.ChartPreset) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.ChartPreset, *pgconn.PgError)
}

type ChartPresetService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewChartPresetService(db *gorm.DB, c *config.Configuration) ChartPresetServiceInterface {
	return &ChartPresetService{db: db, config: c}
}

// PUBLIC FUNCTIONS

func (service *ChartPresetService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.ChartPreset, *pgconn.PgError) {
	var presets []*model.ChartPreset
	result := service.db.Where("thing_id = ?", thingId).Find(&presets)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return presets, nil
}

func (service *ChartPresetService) Create(ctx context.Context, chartPreset *model.ChartPreset) *pgconn.PgError {
	for _, chart := range chartPreset.Charts {
		chart.ChartPresetId = chartPreset.Id
	}
	result := service.db.Create(&chartPreset)
	return utils.GetPostgresError(result.Error)
}

func (service *ChartPresetService) Update(ctx context.Context, updatedChartPreset *model.ChartPreset) *pgconn.PgError {
	result := service.db.Updates(updatedChartPreset)
	return utils.GetPostgresError(result.Error)
}

func (service *ChartPresetService) Delete(ctx context.Context, chartPresetId uuid.UUID) *pgconn.PgError {
	preset := model.ChartPreset{Base: model.Base{Id: chartPresetId}}
	result := service.db.Delete(&preset)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *ChartPresetService) FindById(ctx context.Context, chartPresetId uuid.UUID) (*model.ChartPreset, *pgconn.PgError) { // Should this return a copy or pointer?
	var preset *model.ChartPreset
	result := service.db.Where("id = ?", chartPresetId).First(&preset)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return preset, nil
}
