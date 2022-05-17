package services

import (
	"context"
	"database-ms/app/models"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"
	"errors"
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
		result, err := service.SensorCollection(ctx).InsertOne(ctx, sensor)
		if err == nil {
			sensor.ID = (result.InsertedID).(primitive.ObjectID)
		}
		return err
	}
}

func (service *SensorService) FindByThingId(ctx context.Context, thingId string) ([]*model.Sensor, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	var sensors []*model.Sensor
	cursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	if err = cursor.All(ctx, &sensors); err != nil {
		return nil, err
	}
	if sensors == nil {
		sensors = []*model.Sensor{}
	}
	return sensors, nil
}

func (service *SensorService) FindBySensorId(ctx context.Context, sensorId string) (*model.Sensor, error) {
	bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return nil, err
	}
	var sensor model.Sensor
	if err = service.SensorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonSensorId}).Decode(&sensor); err != nil {
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
	cursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}, {"lastUpdate", bson.D{{"$gt", lastUpdate}}}})
	if err = cursor.All(ctx, &sensors); err != nil {
		return nil, err
	}
	if sensors == nil {
		sensors = []*model.Sensor{}
	}
	return sensors, nil
}

func (service *SensorService) Update(ctx context.Context, updatedSensor *model.Sensor) error {
	sensor, err := service.FindBySensorId(ctx, updatedSensor.ID.Hex())
	if err == nil {
		updatedSensor.SmallId = sensor.SmallId
		updatedSensor.LastUpdate = utils.CurrentTimeInMilli() 
		_, err = service.SensorCollection(ctx).ReplaceOne(ctx, bson.M{"_id": updatedSensor.ID}, updatedSensor)
		return err
	} else {
		return err
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

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		if _, err := db.Collection("Sensor").DeleteOne(ctx, bson.M{"_id": bsonSensorId}); err != nil {
			return nil, err
		}	

		// Put the two below into a function

		// Remove sensor ID from associated raw data presets
		cursor, err := db.Collection("RawDataPreset").Find(ctx, bson.M{"sensorIds": bson.M{"$in": []primitive.ObjectID{bsonSensorId}}})
		if err == nil {
			var rawDataPresets []*model.RawDataPreset
			if err = cursor.All(ctx, &rawDataPresets); err == nil {
				for _, rawDataPreset := range rawDataPresets {
					var sensorIds []primitive.ObjectID
					for _, sensorId := range rawDataPreset.SensorIds {
						if sensorId != bsonSensorId {
							sensorIds = append(sensorIds, sensorId)
						}
					}
					rawDataPreset.SensorIds = sensorIds
					if len(sensorIds) == 0 {
						if _, err := db.Collection("RawDataPreset").DeleteOne(ctx, bson.M{"_id": rawDataPreset.ID}); err != nil {
							return nil, err
						}
					} else {
						if _, err := db.Collection("RawDataPreset").ReplaceOne(ctx, bson.M{"_id": rawDataPreset.ID}, rawDataPreset); err != nil {
							return nil, err
						}
					}
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}

		// Remove sensor ID for associated charts
		chartPresetIdsWithDeletions := []primitive.ObjectID{}
		cursor, err = db.Collection("Chart").Find(ctx, bson.M{"sensorIds": bson.M{"$in": []primitive.ObjectID{bsonSensorId}}})
		if err == nil {
			var charts []*model.Chart
			if err = cursor.All(ctx, &charts); err == nil {
				for _, chart := range charts {
					var sensorIds []primitive.ObjectID
					for _, sensorId := range chart.SensorIds {
						if sensorId != bsonSensorId {
							sensorIds = append(sensorIds, sensorId)
						}
					}
					chart.SensorIds = sensorIds
					if len(sensorIds) == 0 {
						if _, err := db.Collection("Chart").DeleteOne(ctx, bson.M{"_id": chart.ID}); err != nil {
							return nil, err
						}
						chartPresetIdsWithDeletions = append(chartPresetIdsWithDeletions, chart.ChartPresetID)
					} else {
						if _, err := db.Collection("Chart").ReplaceOne(ctx, bson.M{"_id": chart.ID}, chart); err != nil {
							return nil, err
						}
					}
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}

		// Delete any chart presets if they have no charts
		for _, chartPresetId := range chartPresetIdsWithDeletions {
			cursor, err := db.Collection("Chart").Find(ctx, bson.M{"chartPresetId": chartPresetId})
			var charts []*models.Chart
			if err = cursor.All(ctx, &charts); err == nil {
				if len(charts) == 0 {
					if _, err := db.Collection("ChartPreset").DeleteOne(ctx, bson.M{"_id": chartPresetId}); err != nil {
						return nil, err
					}
				}
			} else {
				return nil, err
			}		
		}
		
		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *SensorService) IsSensorUnique(ctx context.Context, newSensor *model.Sensor) bool {
	// TODO: Do with FindOne query rather than fetching everything
	sensors, err := service.FindByThingId(ctx, newSensor.ThingID.Hex())
	if err == nil {
		for _, sensor := range sensors {
			if (newSensor.Name == sensor.Name || newSensor.CanId == sensor.CanId) && newSensor.ID != sensor.ID {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// ============== Service Helper Method(s) ================

func (service *SensorService) FindAvailableSmallId(thingId primitive.ObjectID, ctx context.Context) (int, error) {
	opts := options.Find().SetProjection(bson.D{{"smallId", 1}, {"_id", 0}})
	filterCursor, err := service.SensorCollection(ctx).Find(ctx, bson.D{{"thingId", thingId}}, opts)
	if err != nil {
		return -1, err
	}

	type SmallId struct {
		SmallId int
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

func (service *SensorService) SensorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Sensor")
}
