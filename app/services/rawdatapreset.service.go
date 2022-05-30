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

type RawDataPresetServiceInterface interface {
	// Public
	Create(context.Context, *model.RawDataPreset) *pgconn.PgError
	FindByThingId(context.Context, uuid.UUID) ([]*model.RawDataPreset, *pgconn.PgError)
	Update(context.Context, *model.RawDataPreset) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError
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
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get all the presets associated with the given thing
		result := db.Where("thing_id = ?", thingId).Find(&presets)
		if result.Error != nil {
			return result.Error
		}

		// Get the ids of the relationship with sensor
		for _, preset := range presets {
			var presetSensors []*model.RawDataPresetSensor
			preset.SensorIds = []uuid.UUID{}
			result = db.Table(model.TableNameRawdataPresetSensor).Where("rawdatapreset_id = ?", preset.Id).Find(&presetSensors)
			if result.Error != nil {
				return result.Error
			}
			for _, presetSensor := range presetSensors {
				preset.SensorIds = append(preset.SensorIds, presetSensor.SensorId)
			}
		}
		return result.Error
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return presets, nil
}

func (service *RawDataPresetService) Create(ctx context.Context, rawDataPreset *model.RawDataPreset) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Create the preset
		result := db.Create(&rawDataPreset)
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of preset sensors
		var rawDataPresetSensors []model.RawDataPresetSensor
		for _, sensorId := range rawDataPreset.SensorIds {
			rawDataPresetSensors = append(rawDataPresetSensors, model.RawDataPresetSensor{
				RawDataPresetId: rawDataPreset.Id,
				SensorId:        sensorId,
			})
		}

		// Insert empty sensorIds
		if len(rawDataPreset.SensorIds) == 0 {
			rawDataPreset.SensorIds = []uuid.UUID{}
		}

		// Batch insert preset sensor
		result = db.Table(model.TableNameRawdataPresetSensor).CreateInBatches(rawDataPresetSensors, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
}

func (service *RawDataPresetService) Update(ctx context.Context, updatedRawDataPreset *model.RawDataPreset) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Save the updated preset
		result := db.Updates(updatedRawDataPreset)
		if result.Error != nil {
			return result.Error
		}

		// Delete all of the associated preset-sensors
		result = db.Table(model.TableNameRawdataPresetSensor).Where("rawdatapreset_id = ?", updatedRawDataPreset.Id).Delete(&model.RawDataPresetSensor{})
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of preset-sensors
		presetSensors := []model.RawDataPresetSensor{}
		for _, sensorId := range updatedRawDataPreset.SensorIds {
			presetSensor := model.RawDataPresetSensor{}
			presetSensor.RawDataPresetId = updatedRawDataPreset.Id
			presetSensor.SensorId = sensorId
			presetSensors = append(presetSensors, presetSensor)
		}

		// Batch insert preset-sensors
		result = db.Table(model.TableNameRawdataPresetSensor).CreateInBatches(presetSensors, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
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
		return nil, utils.GetPostgresError(result.Error)
	}
	return preset, nil
}
