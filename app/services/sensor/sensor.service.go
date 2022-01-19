package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/utils"
	"errors"
	"sort"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SensorServiceInterface interface {
	Create(context.Context, *model.Sensor) error
	FindByThingId(context.Context, string) ([]*model.Sensor, error)
	FindBySensorId(context.Context, string) (*model.Sensor, error)
	FindUpdatedSensor(context.Context, string, int64) ([]*model.Sensor, error)
	Update(context.Context, string, *model.SensorUpdate) error
	Delete(context.Context, string) error
}

type SensorService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewSensorService(db *mgo.Session, c *config.Configuration) SensorServiceInterface {
	return &SensorService{config: c, db: db}
}

func (service *SensorService) Create(ctx context.Context, sensor *model.Sensor) error {
	newSmallId, err := service.findAvailableSmallId(sensor.ThingID)
	if err != nil {
		return err
	}
	sensor.SmallId = &newSmallId
	sensor.ID = bson.NewObjectId()
	sensor.LastUpdate = utils.CurrentTimeInMilli()

	// TODO add new sensor to Thing SensorId list

	return service.addSensor(ctx, sensor)
}

func (service *SensorService) FindByThingId(ctx context.Context, thingId string) ([]*model.Sensor, error) {

	return service.getSensors(ctx, bson.M{"thingId": bson.ObjectIdHex(thingId)})

}

func (service *SensorService) FindBySensorId(ctx context.Context, sensorId string) (*model.Sensor, error) {

	return service.getSensor(ctx, bson.M{"thingId": bson.ObjectIdHex(sensorId)})

}

func (service *SensorService) FindUpdatedSensor(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Sensor, error) {

	return service.getSensors(ctx, bson.M{
		"thingId": bson.ObjectIdHex(thingId),
		"lastUpdate": bson.M{
			"$gt": lastUpdate,
		},
	})

}

func (service *SensorService) Update(ctx context.Context, sensorId string, sensor *model.SensorUpdate) error {
	query := bson.M{"_id": bson.ObjectIdHex(sensorId)}
	CustomBson := &utils.CustomBson{}
	change, err := CustomBson.Set(sensor)
	if err != nil {
		return err
	}

	return service.collection().Update(query, change)
}

func (service *SensorService) Delete(ctx context.Context, sensorId string) error {

	// TODO Delete sensorId from Thing SensorId list

	return service.collection().RemoveId(bson.ObjectIdHex(sensorId))

}

// ============== Service Helper Method(s) ================

func (service *SensorService) findAvailableSmallId(thingId bson.ObjectId) (int, error) {
	var result []*model.Sensor
	err := service.collection().Find(bson.M{"thingId": thingId}).Select(bson.M{"smallId": 1}).All(&result)
	if err != nil {
		return -1, err
	}

	var smallIds []int
	for _, record := range result {
		smallIds = append(smallIds, *record.SmallId)
	}
	smallIds = utils.Unique(smallIds)
	sort.Ints(smallIds)

	availableSmallId := 0
	for _, smallId := range smallIds {
		if smallId != availableSmallId {
			break
		}
		availableSmallId++
	}

	if availableSmallId < 256 {
		return availableSmallId, nil
	} else {
		return 0, errors.New("no available smallIds")
	}
}

// ============== Common DB Operations ===================

func (service *SensorService) addSensor(ctx context.Context, sensor *model.Sensor) error {
	return service.collection().Insert(sensor)
}

func (service *SensorService) getSensor(ctx context.Context, query interface{}) (*model.Sensor, error) {
	var sensor model.Sensor
	err := service.collection().Find(query).One(&sensor)
	return &sensor, err
}

func (service *SensorService) getSensors(ctx context.Context, query interface{}) ([]*model.Sensor, error) {
	var sensors []*model.Sensor
	err := service.collection().Find(query).All(&sensors)
	return sensors, err
}

func (service *SensorService) collection() *mgo.Collection {
	return service.db.DB(service.config.MongoDbName).C("Sensor")
}
