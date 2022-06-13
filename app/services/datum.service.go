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

type DatumServiceInterface interface {
	// Public
	FindBySessionIdAndSensorId(context.Context, uuid.UUID, uuid.UUID) ([]*SensorData, *pgconn.PgError)
	CreateMany(context.Context, []*model.Datum) *pgconn.PgError
}

type DatumService struct {
	db     *gorm.DB
	config *config.Configuration
}

type SensorData struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

func NewDatumService(db *gorm.DB, c *config.Configuration) DatumServiceInterface {
	return &DatumService{config: c, db: db}
}

func (service *DatumService) FindBySessionIdAndSensorId(ctx context.Context, sessionId uuid.UUID, sensorId uuid.UUID) ([]*SensorData, *pgconn.PgError) {
	var data []*model.Datum
	result := service.db.Order("timestamp asc").Find(&data, "session_id = ? AND sensor_id = ?", sessionId, sensorId)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	print(len(data))
	cleanData := []*SensorData{}
	for _, datum := range data {
		cleanData = append(cleanData, &SensorData{X: datum.Timestamp, Y: datum.Value})
	}
	return cleanData, nil
}

func (service *DatumService) CreateMany(ctx context.Context, datumArray []*model.Datum) *pgconn.PgError {
	result := service.db.CreateInBatches(datumArray, 100)
	return utils.GetPostgresError(result.Error)
}
