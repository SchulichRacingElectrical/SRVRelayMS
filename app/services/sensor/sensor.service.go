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
