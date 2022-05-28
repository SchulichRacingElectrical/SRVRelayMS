package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SensorServiceInterface interface {
	Create(context.Context, *model.Sensor) error
	FindByThingId(context.Context, uuid.UUID) ([]*model.Sensor, error)
	FindBySensorId(context.Context, uuid.UUID) (*model.Sensor, error)
	FindUpdatedSensors(context.Context, uuid.UUID, int64) ([]*model.Sensor, error)
	Update(context.Context, *model.Sensor) error
	Delete(context.Context, uuid.UUID) error
	IsSensorUnique(context.Context, *model.Sensor) bool
}

type SensorService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewSensorService(db *gorm.DB, c *config.Configuration) SensorServiceInterface {
	return &SensorService{config: c, db: db}
}

func (service *SensorService) Create(ctx context.Context, sensor *model.Sensor) error {
	// newSmallId, err := service.FindAvailableSmallId(sensor.ThingID, ctx)
	// if err != nil {
	// 	return err
	// } else {
	// 	sensor.SmallId = &newSmallId
	// 	sensor.LastUpdate = utils.CurrentTimeInMilli()
	// 	result, err := service.SensorCollection(ctx).InsertOne(ctx, sensor)
	// 	if err == nil {
	// 		sensor.ID = (result.InsertedID).(primitive.ObjectID)
	// 	}
	// 	return err
	// }
	return nil
}

func (service *SensorService) FindByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Sensor, error) {
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }
	// var sensors []*model.Sensor
	// cursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	// if err = cursor.All(ctx, &sensors); err != nil {
	// 	return nil, err
	// }
	// if sensors == nil {
	// 	sensors = []*model.Sensor{}
	// }
	// return sensors, nil
	return nil, nil
}

func (service *SensorService) FindBySensorId(ctx context.Context, sensorId uuid.UUID) (*model.Sensor, error) {
	// bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	// if err != nil {
	// 	return nil, err
	// }
	// var sensor model.Sensor
	// if err = service.SensorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonSensorId}).Decode(&sensor); err != nil {
	// 	return nil, err
	// }
	// return &sensor, nil
	return nil, nil
}

func (service *SensorService) FindUpdatedSensors(ctx context.Context, thingId uuid.UUID, lastUpdate int64) ([]*model.Sensor, error) {
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }
	// var sensors []*model.Sensor
	// cursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}, {"lastUpdate", bson.D{{"$gt", lastUpdate}}}})
	// if err = cursor.All(ctx, &sensors); err != nil {
	// 	return nil, err
	// }
	// if sensors == nil {
	// 	sensors = []*model.Sensor{}
	// }
	// return sensors, nil
	return nil, nil
}

func (service *SensorService) Update(ctx context.Context, updatedSensor *model.Sensor) error {
	// sensor, err := service.FindBySensorId(ctx, updatedSensor.ID.Hex())
	// if err == nil {
	// 	updatedSensor.SmallId = sensor.SmallId
	// 	updatedSensor.LastUpdate = utils.CurrentTimeInMilli()
	// 	_, err = service.SensorCollection(ctx).ReplaceOne(ctx, bson.M{"_id": updatedSensor.ID}, updatedSensor)
	// 	return err
	// } else {
	// 	return err
	// }
	return nil
}

func (service *SensorService) Delete(ctx context.Context, sensorId uuid.UUID) error {
	// bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	// if err != nil {
	// 	return err
	// } else {
	// 	_, err := service.SensorCollection(ctx).DeleteOne(ctx, bson.M{"_id": bsonSensorId})
	// 	return err
	// }
	return nil
}

func (service *SensorService) IsSensorUnique(ctx context.Context, newSensor *model.Sensor) bool {
	// sensors, err := service.FindByThingId(ctx, newSensor.ThingID.Hex())
	// if err == nil {
	// 	for _, sensor := range sensors {
	// 		if (newSensor.Name == sensor.Name || newSensor.CanId == sensor.CanId) && newSensor.ID != sensor.ID {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// } else {
	// 	return false
	// }
	return false
}

// ============== Service Helper Method(s) ================

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
