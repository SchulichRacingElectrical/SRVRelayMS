package sensor

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
	newSmallId, err := service.findAvailableSmallId(sensor.ThingID, ctx)
	if err != nil {
		return err
	}

	sensor.SmallId = &newSmallId
	sensor.ID = primitive.NewObjectID()
	sensor.LastUpdate = utils.CurrentTimeInMilli()

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Check if thing exists
		thingCollection := db.Collection("Thing")
		if count, err := thingCollection.CountDocuments(sessCtx, bson.M{"_id": sensor.ThingID}); err != nil || count < 1 {
			return nil, errors.New("could not find thing")
		}

		// Insert new sensor
		fmt.Println("Inserting new sensor...")
		sensorCollection := db.Collection("Sensor")
		_, err := sensorCollection.InsertOne(ctx, sensor)
		if err != nil {
			return nil, err
		}

		// Add sensor id to thing sensor list
		fmt.Println("Adding sensor to thing...")
		updpate := bson.M{"$push": bson.M{"sensors": sensor.ID}}
		if _, err := thingCollection.UpdateByID(ctx, sensor.ThingID, updpate); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
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

func (service *SensorService) FindUpdatedSensor(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Sensor, error) {
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

func (service *SensorService) Update(ctx context.Context, sensorId string, updates *model.SensorUpdate) error {
	bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return err
	}
	_, err = service.sensorCollection(ctx).UpdateOne(ctx, bson.M{"_id": bsonSensorId}, bson.M{"$set": updates})
	return err
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

// ============== Service Helper Method(s) ================
type SmallId struct {
	SmallId int
}

func (service *SensorService) findAvailableSmallId(thingId primitive.ObjectID, ctx context.Context) (int, error) {
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
