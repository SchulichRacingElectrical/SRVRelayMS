package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/app/utils"
	"database-ms/config"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type RawDataPresetServiceInterface interface {
	// Public
	FindByThingId(context.Context, uuid.UUID) ([]*model.RawDataPreset, *pgconn.PgError)
	Create(context.Context, *model.RawDataPreset) *pgconn.PgError
	Update(context.Context, *model.RawDataPreset) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.RawDataPreset, *pgconn.PgError)
}

type RawDataPresetService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewRawDataPresetService(db *gorm.DB, c *config.Configuration) RawDataPresetServiceInterface {
	return &RawDataPresetService{db: db, config: c}
}

// PUBLIC FUNCTIONS

func (service *RawDataPresetService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.RawDataPreset, *pgconn.PgError) {
	var presets []*model.RawDataPreset
	result := service.db.Where("thing_id = ?", thingId).Find(&presets)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return presets, nil
}

func (service *RawDataPresetService) Create(ctx context.Context, rawDataPreset *model.RawDataPreset) *pgconn.PgError {
	result := service.db.Create(&rawDataPreset)
	return utils.GetPostgresError(result.Error)
}

func (service *RawDataPresetService) Update(ctx context.Context, updatedRawDataPreset *model.RawDataPreset) *pgconn.PgError {
	result := service.db.Updates(updatedRawDataPreset)
	return utils.GetPostgresError(result.Error)
}

func (service *RawDataPresetService) Delete(ctx context.Context, rawDataPresetId uuid.UUID) *pgconn.PgError {
	preset := model.RawDataPreset{Base: model.Base{Id: rawDataPresetId}}
	result := service.db.Delete(&preset)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *RawDataPresetService) FindById(ctx context.Context, rawDataPresetId uuid.UUID) (*model.RawDataPreset, *pgconn.PgError) {
	var preset *model.RawDataPreset
	result := service.db.Where("id = ?", rawDataPresetId).First(&preset)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return preset, nil
}
