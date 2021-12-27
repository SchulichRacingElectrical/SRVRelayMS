package user

import (
	"context"
	model "database-ms/app/models"
	repository "database-ms/app/repositories/sensor"
	"database-ms/config"
	"database-ms/utils"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorServiceInterface interface {
	Create(context.Context, *model.Sensor) error
	FindByThingId(context.Context, string) ([]*model.Sensor, error)
	FindById(context.Context, string) (*model.Sensor, error)
	FindByThingIdAndSid(context.Context, string, int) (*model.Sensor, error)
	FindByThingIdAndLastUpdate(context.Context, string, int64) ([]*model.Sensor, error)
	Update(context.Context, string, *model.SensorUpdate) error
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

func (service *SensorService) Update(ctx context.Context, id string, sensor *model.SensorUpdate) error {
	// TODO handle id and thingId + sid
	query := bson.M{"_id": bson.ObjectIdHex(id)}
	CustomBson := &utils.CustomBson{}
	change, err := CustomBson.Set(sensor)
	if err != nil {
		return err
	}
	return service.repository.Update(ctx, query, change)
}
