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

type DatumServiceInterface interface {
	// Public
	FindBySessionIdAndSensorId(context.Context, uuid.UUID, uuid.UUID) ([]SensorData, *pgconn.PgError)
	CreateMany(context.Context, []*model.Datum) *pgconn.PgError
}

type DatumService struct {
	db     *gorm.DB
	config *config.Configuration
}

type SensorData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func NewDatumService(db *gorm.DB, c *config.Configuration) DatumServiceInterface {
	return &DatumService{config: c, db: db}
}

func (service *DatumService) FindBySessionIdAndSensorId(ctx context.Context, sessionId uuid.UUID, sensorId uuid.UUID) ([]SensorData, *pgconn.PgError) {
	var data []*model.Datum
	result := service.db.Where("session_id = ? AND sensor_id = ?", sessionId, sensorId).Find(&data)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	// TODO: Sort the data by timestamp
	// TODO: Create array of SensorData
	// TODO: Return the data
	return nil, nil
}

func (service *DatumService) CreateMany(ctx context.Context, datumArray []*model.Datum) *pgconn.PgError {
	result := service.db.CreateInBatches(datumArray, 100)
	return utils.GetPostgresError(result.Error)
}
