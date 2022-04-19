package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"
	"errors"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
)

type SensorServiceInterface interface {
	Create(context.Context, *model.Sensor) error
	FindByThingId(context.Context, string) ([]*model.Sensor, error)
	FindBySensorId(context.Context, string) (*model.Sensor, error)
	FindUpdatedSensors(context.Context, string, int64) ([]*model.Sensor, error)
	Update(context.Context, *model.Sensor) error
	Delete(context.Context, string) error
	IsSensorUnique(context.Context, *model.Sensor) bool
}

type SensorService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewSensorService(db *mgo.Session, c *config.Configuration) SensorServiceInterface {
	return &SensorService{config: c, db: db}
}

func (service *SensorService) Create(ctx context.Context, sensor *model.Sensor) error {
	newSmallId, err := service.FindAvailableSmallId(sensor.ThingID, ctx)
	if err != nil {
		return err
	} else {
		sensor.SmallId = &newSmallId
		sensor.LastUpdate = utils.CurrentTimeInMilli()
		result, err := service.sensorCollection(ctx).InsertOne(ctx, sensor)
		sensor.ID = (result.InsertedID).(primitive.ObjectID)
		return err
	}
}

func (service *SensorService) FindByThingId(ctx context.Context, thingId string) ([]*model.Sensor, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	var sensors []*model.Sensor
	cursor, err := service.sensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	if err = cursor.All(ctx, &sensors); err != nil {
		return nil, err
	}
	return sensors, nil
}

func (service *SensorService) FindBySensorId(ctx context.Context, sensorId string) (*model.Sensor, error) {
	bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return nil, err
	}
	var sensor model.Sensor
	if err = service.sensorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonSensorId}).Decode(&sensor); err != nil {
		return nil, err
	}
	return &sensor, nil
}

func (service *SensorService) FindUpdatedSensors(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Sensor, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	var sensors []*model.Sensor
	cursor, err := service.sensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}, {"lastUpdate", bson.D{{"$gt", lastUpdate}}}})
	if err = cursor.All(ctx, &sensors); err != nil {
		return nil, err
	}
	return sensors, nil
}

func (service *SensorService) Update(ctx context.Context, updatedSensor *model.Sensor) error {
	sensor, err := service.FindBySensorId(ctx, updatedSensor.ID.String())
	if err != nil {
		return err
	} else {
		updatedSensor.SmallId = sensor.SmallId
		if service.IsSensorUnique(ctx, updatedSensor) {
			_, err = service.sensorCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedSensor.ID}, bson.M{"$set": updatedSensor})
			return err
		} else {
			return errors.New("Sensor name and/or CAN ID must remain unique.")
		}
	}
}

func (service *SensorService) Delete(ctx context.Context, sensorId string) error {
	bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		sensorCollection := db.Collection("Sensor")
		// Get sensor
		fmt.Println("Getting sensor...")
		var sensor model.Sensor
		if err = service.sensorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonSensorId}).Decode(&sensor); err != nil {
			return nil, err
		}

		// Delete sensor from thing sensors list
		thingCollection := db.Collection("Thing")
		fmt.Println("Removing sensor from thing...")
		updpate := bson.M{"$pull": bson.M{"sensors": bsonSensorId}}
		if _, err := thingCollection.UpdateByID(ctx, sensor.ThingID, updpate); err != nil {
			return nil, err
		}

		// Delete sensor
		fmt.Println("Deleting sensor...")
		_, err := sensorCollection.DeleteOne(ctx, bson.M{"_id": bsonSensorId})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *SensorService) IsSensorUnique(ctx context.Context, newSensor *model.Sensor) bool {
	sensors, err := service.FindByThingId(ctx, newSensor.ThingID.String())
	if err != nil {
		return false
	} else {
		for _, sensor := range sensors {
			if newSensor.Name == sensor.Name || newSensor.CanId == sensor.CanId {
				return false
			}
		}
		return true
	}
}

// ============== Service Helper Method(s) ================
type SmallId struct {
	SmallId int
}

func (service *SensorService) FindAvailableSmallId(thingId primitive.ObjectID, ctx context.Context) (int, error) {
	opts := options.Find().SetProjection(bson.D{{"smallId", 1}, {"_id", 0}})
	filterCursor, err := service.sensorCollection(ctx).Find(ctx, bson.D{{"thingId", thingId}}, opts)
	if err != nil {
		return -1, err
	}

	var results []SmallId
	if err = filterCursor.All(ctx, &results); err != nil {
		return -1, err
	}

	var smallIds []int
	for _, record := range results {
		smallIds = append(smallIds, record.SmallId)
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

func (service *SensorService) sensorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}

	return dbClient.Database(service.config.MongoDbName).Collection("Sensor")
}
