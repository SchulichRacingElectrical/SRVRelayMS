package services

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type DatumServiceInterface interface {
	Create(context.Context, *models.Datum) error
	CreateMany(context.Context, []*models.Datum) error
}

type DatumService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewDatumService(db *mgo.Session, c *config.Configuration) DatumServiceInterface {
	return &DatumService{config: c, db: db}
}

func (service *DatumService) Create(ctx context.Context, datum *models.Datum) error {
	result, err := service.DatumCollection(ctx).InsertOne(ctx, datum)
	if err == nil {
		datum.ID = (result.InsertedID).(primitive.ObjectID)
	}
	return err
}

func (service *DatumService) CreateMany(ctx context.Context, datumArray []*models.Datum) error {
	docs := make([]interface{}, len(datumArray))
	for i, datum := range datumArray {
		docs[i] = datum
	}
	_, err := service.DatumCollection(ctx).InsertMany(ctx, docs)
	return err
}

// ============== Common DB Operations ===================

func (service *DatumService) DatumCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Datum")
}
