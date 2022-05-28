package services

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type DatumServiceInterface interface {
	Create(context.Context, *models.Datum) error
	CreateMany(context.Context, []*models.Datum) error
	FindBySessionIdAndSensorId(context.Context, string, string) ([]*models.Datum, error)
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

func (service *DatumService) FindBySessionIdAndSensorId(ctx context.Context, sessionId string, sensorId string) ([]*models.Datum, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	var datumArray []*models.Datum
	bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return nil, err
	}
	bsonSensorId, err := primitive.ObjectIDFromHex(sensorId)
	if err != nil {
		return nil, err
	}

	cursor, err := database.Collection("Datum").Find(ctx, bson.M{"sessionId": bsonSessionId, "sensorId": bsonSensorId})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &datumArray); err != nil {
		return nil, err
	}

	return datumArray, nil
}

// ============== Common DB Operations ===================

func (service *DatumService) DatumCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Datum")
}
