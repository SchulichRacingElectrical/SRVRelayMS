package services

import (
	"context"
	"database-ms/app/models"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type RawDataPresetServiceInterface interface {
	Create(context.Context, *model.RawDataPreset) error
	FindByThingId(context.Context, string) ([]*model.RawDataPreset, error)
	Update(context.Context, *model.RawDataPreset) error
	Delete(context.Context, string) error
	FindById(context.Context, string) (*model.RawDataPreset, error)
	IsPresetUnique(context.Context, *model.RawDataPreset) bool
	IsPresetValid(context.Context, *model.RawDataPreset) bool
}

type RawDataPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewRawDataPresetService(db *mgo.Session, c *config.Configuration) RawDataPresetServiceInterface {
	return &RawDataPresetService{db: db, config: c}
}

func (service *RawDataPresetService) Create(ctx context.Context, rawDataPreset *model.RawDataPreset) error {
	// Remove duplicate sensor Ids from the preset
	sensorIdMap := make(map[primitive.ObjectID]int)
	for _, sensorId := range rawDataPreset.SensorIds {
		sensorIdMap[sensorId] = 0
	}
	rawDataPreset.SensorIds = []primitive.ObjectID{}
	for id, _ := range sensorIdMap {
		rawDataPreset.SensorIds = append(rawDataPreset.SensorIds, id)
	}
	result, err := service.RawDataPresetCollection(ctx).InsertOne(ctx, rawDataPreset)
	if err == nil {
		rawDataPreset.ID = (result.InsertedID).(primitive.ObjectID)
	}
	return err
}

func (service *RawDataPresetService) FindByThingId(ctx context.Context, thingId string) ([]*model.RawDataPreset, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	var rawDataPresets []*models.RawDataPreset
	cursor, err := service.RawDataPresetCollection(ctx).Find(ctx, bson.D{{"thingId", bsonThingId}})
	if err = cursor.All(ctx, &rawDataPresets); err != nil {
		return nil, err
	}
	if rawDataPresets == nil {
		rawDataPresets = []*models.RawDataPreset{}
	}
	return rawDataPresets, nil
}

func (service *RawDataPresetService) Update(ctx context.Context, updatedRawDataPreset *model.RawDataPreset) error {
	// Remove duplicate sensor Ids from the preset
	sensorIdMap := make(map[primitive.ObjectID]int)
	for _, sensorId := range updatedRawDataPreset.SensorIds {
		sensorIdMap[sensorId] = 0
	}
	updatedRawDataPreset.SensorIds = []primitive.ObjectID{}
	for id, _ := range sensorIdMap {
		updatedRawDataPreset.SensorIds = append(updatedRawDataPreset.SensorIds, id)
	}
	_, err := service.RawDataPresetCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedRawDataPreset.ID}, bson.M{"$set": updatedRawDataPreset})
	return err
}

func (service *RawDataPresetService) Delete(ctx context.Context, rawDataPresetId string) error {
	bsonRawDataPresetId, err := primitive.ObjectIDFromHex(rawDataPresetId)
	if err == nil {
		_, err := service.RawDataPresetCollection(ctx).DeleteOne(ctx, bson.M{"_id": bsonRawDataPresetId})
		return err
	} else {
		return err
	}
}

func (service *RawDataPresetService) FindById(ctx context.Context, rawDataPresetId string) (*model.RawDataPreset, error) {
	bsonRawDataPresetId, err := primitive.ObjectIDFromHex(rawDataPresetId)
	if err != nil {
		return nil, err
	}
	var rawDataPreset models.RawDataPreset
	if err = service.RawDataPresetCollection(ctx).FindOne(ctx, bson.M{"_id": bsonRawDataPresetId}).Decode(&rawDataPreset); err != nil {
		return nil, err
	}
	return &rawDataPreset, nil
}

func (service *RawDataPresetService) IsPresetUnique(ctx context.Context, newRawDataPreset *model.RawDataPreset) bool {
	// TODO: Do with FindOne query rather than fetching everything
	rawDataPresets, err := service.FindByThingId(ctx, newRawDataPreset.ThingId.Hex())
	if err == nil {
		for _, rawDataPreset := range rawDataPresets {
			if newRawDataPreset.Name == rawDataPreset.Name && newRawDataPreset.ID != rawDataPreset.ID {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *RawDataPresetService) IsPresetValid(ctx context.Context, rawDataPreset *model.RawDataPreset) bool {
	_, err := service.SensorCollection(ctx).Find(ctx, bson.M{"_id": bson.M{"$in": rawDataPreset.SensorIds }})
	return err == nil
}

func (service *RawDataPresetService) SensorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Sensor")
}

func (service *RawDataPresetService) RawDataPresetCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("RawDataPreset")
}