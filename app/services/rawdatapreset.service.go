package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

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
	IsRawDataPresetUnique(context.Context, *model.RawDataPreset) bool
	DoPresetSensorsExist(context.Context, *model.RawDataPreset) bool
}

type RawDataPresetService struct {
	db			*mgo.Session
	config	*config.Configuration
}

func NewRawDataPresetService(db *mgo.Session, c *config.Configuration) RawDataPresetServiceInterface {
	return &RawDataPresetService{db: db, config: c}
}

func (service *RawDataPresetService) Create(ctx context.Context, rawDataPreset *model.RawDataPreset) error {
	result, err := service.RawDataPresetCollection(ctx).InsertOne(ctx, rawDataPreset)
	if err == nil {
		rawDataPreset.ID = (result.InsertedID).(primitive.ObjectID)
	}
	return err
}

func (service *RawDataPresetService) FindByThingId(ctx context.Context, thingId string) ([]*model.RawDataPreset, error) {
	return nil, nil
}

func (service *RawDataPresetService) Update(ctx context.Context, updatedRawDataPreset *model.RawDataPreset) error {
	return nil
}

func (service *RawDataPresetService) Delete(ctx context.Context, rawDataPresetId string) error {
	return nil
}

func (service *RawDataPresetService) FindById(ctx context.Context, rawDataPresetId string) (*model.RawDataPreset, error) {
	return nil, nil
}

func (service *RawDataPresetService) IsRawDataPresetUnique(ctx context.Context, rawDataPreset *model.RawDataPreset) bool {
	return false
}

func (service *RawDataPresetService) DoPresetSensorsExist(ctx context.Context, rawDataPreset *model.RawDataPreset) bool {
	return false
}

func (service *RawDataPresetService) RawDataPresetCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("RawDataPreset")
}