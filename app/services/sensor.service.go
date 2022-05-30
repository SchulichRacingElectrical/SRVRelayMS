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

type SensorServiceInterface interface {
	// Public
	FindByThingId(context.Context, uuid.UUID) ([]*model.Sensor, *pgconn.PgError)
	FindUpdatedSensors(context.Context, uuid.UUID, int64) ([]*model.Sensor, *pgconn.PgError)
	Create(context.Context, *model.Sensor) *pgconn.PgError
	Update(context.Context, *model.Sensor) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindBySensorId(context.Context, uuid.UUID) (*model.Sensor, *pgconn.PgError)
	FindAvailableSmallId(uuid.UUID, context.Context) (int, error)
}

type SensorService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewSensorService(db *gorm.DB, c *config.Configuration) SensorServiceInterface {
	return &SensorService{config: c, db: db}
}

// PUBLIC FUNCTIONS

func (service *SensorService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Sensor, *pgconn.PgError) {
	var sensors []*model.Sensor
	result := service.db.Where("thing_id = ?", thingId).Find(&sensors)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return sensors, nil
}

func (service *SensorService) FindUpdatedSensors(ctx context.Context, thingId uuid.UUID, lastUpdate int64) ([]*model.Sensor, *pgconn.PgError) {
	var sensors []*model.Sensor
	result := service.db.Where("thing_id = ? AND last_update > ?", thingId, lastUpdate).Find(&sensors)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return sensors, nil
}

func (service *SensorService) Create(ctx context.Context, sensor *model.Sensor) *pgconn.PgError {
	// Generate a small id
	newSmallId, err := service.FindAvailableSmallId(sensor.ThingId, ctx)
	if err != nil {
		return utils.GetPostgresError(err) // Check if this works
	}

	// Create the sensor
	sensor.SmallId = int32(newSmallId)
	sensor.LastUpdate = utils.CurrentTimeInMilli()
	result := service.db.Create(&sensor)
	return utils.GetPostgresError(result.Error)
}

func (service *SensorService) Update(ctx context.Context, updatedSensor *model.Sensor) *pgconn.PgError {
	sensor, err := service.FindBySensorId(ctx, updatedSensor.Id)
	if err != nil {
		return err
	}
	updatedSensor.SmallId = sensor.SmallId
	updatedSensor.LastUpdate = utils.CurrentTimeInMilli()
	result := service.db.Save(&updatedSensor)
	return utils.GetPostgresError(result.Error)
}

func (service *SensorService) Delete(ctx context.Context, sensorId uuid.UUID) *pgconn.PgError {
	sensor := model.Sensor{Base: model.Base{Id: sensorId}}
	result := service.db.Delete(&sensor)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *SensorService) FindBySensorId(ctx context.Context, sensorId uuid.UUID) (*model.Sensor, *pgconn.PgError) {
	var sensor *model.Sensor
	result := service.db.Where("id = ?", sensorId).First(&sensor)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return sensor, nil
}

func (service *SensorService) FindAvailableSmallId(thingId uuid.UUID, ctx context.Context) (int, error) {
	// opts := options.Find().SetProjection(bson.D{{"smallId", 1}, {"_id", 0}})
	// filterCursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", thingId}}, opts)
	// if err != nil {
	// 	return -1, err
	// }

	// type SmallId struct {
	// 	SmallId int
	// }
	// var results []SmallId
	// if err = filterCursor.All(ctx, &results); err != nil {
	// 	return -1, err
	// }

	// var smallIds []int
	// for _, record := range results {
	// 	smallIds = append(smallIds, record.SmallId)
	// }

	// smallIds = utils.Unique(smallIds)
	// sort.Ints(smallIds)

	// availableSmallId := 0
	// for _, smallId := range smallIds {
	// 	if smallId != availableSmallId {
	// 		break
	// 	}
	// 	availableSmallId++
	// }

	// if availableSmallId < 256 {
	// 	return availableSmallId, nil
	// } else {
	// 	return 0, errors.New("no available smallIds")
	// }
	return 0, nil
}
