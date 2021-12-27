package user

import (
	"context"
	model "database-ms/app/models"
	repository "database-ms/app/repositories/sensor"
	"database-ms/config"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorServiceInterface interface {
	Create(context.Context, *model.Sensor) error
	IsSensorAlreadyExists(context.Context, bson.ObjectId, int) bool
	FindByThingId(context.Context, string) ([]*model.Sensor, error)
	FindById(context.Context, string) (*model.Sensor, error)
	FindByThingIdAndSid(context.Context, string, int) (*model.Sensor, error)
	FindByThingIdAndLastUpdate(context.Context, string, int64) ([]*model.Sensor, error)
}

type SensorService struct {
	db         *mgo.Session
	repository repository.SensorRepository
	config     *config.Configuration
}

func New(sensorRepo repository.SensorRepository) SensorServiceInterface {
	return &SensorService{repository: sensorRepo}
}

func (service *SensorService) Create(ctx context.Context, sensor *model.Sensor) error {

	return service.repository.Create(ctx, sensor)

}

func (service *SensorService) IsSensorAlreadyExists(ctx context.Context, thingId bson.ObjectId, sid int) bool {

	return service.repository.IsSensorAlreadyExisits(ctx, thingId, sid)

}

func (service *SensorService) FindByThingId(ctx context.Context, thingId string) ([]*model.Sensor, error) {

	return service.repository.FindByThingId(ctx, thingId)

}

func (service *SensorService) FindById(ctx context.Context, id string) (*model.Sensor, error) {

	return service.repository.FindOneById(ctx, id)

}

func (service *SensorService) FindByThingIdAndSid(ctx context.Context, thingId string, sid int) (*model.Sensor, error) {

	return service.repository.FindOneByThingIdAndSid(ctx, thingId, sid)

}

func (service *SensorService) FindByThingIdAndLastUpdate(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Sensor, error) {

	return service.repository.FindByThingIdAndLastUpdate(ctx, thingId, lastUpdate)

}
