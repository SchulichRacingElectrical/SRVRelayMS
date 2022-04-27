package services

import (
	"context"

	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type ThingOperatorServiceInterface interface {
	Create(context.Context, *model.ThingOperator) error
	Delete(context.Context, string, string) error
	IsAssociationUnique(context.Context, *model.ThingOperator) bool
}

type ThingOperatorService struct {
	db			*mgo.Session
	config 	*config.Configuration
}

func NewThingOperatorService(db *mgo.Session, c *config.Configuration) ThingOperatorServiceInterface {
	return &ThingOperatorService{db: db, config: c}
}

func (service *ThingOperatorService) Create(ctx context.Context, thingOperator *model.ThingOperator) error {
	_, err := service.ThingOperatorCollection(ctx).InsertOne(ctx, thingOperator)
	return err
}

func (service *ThingOperatorService) Delete(ctx context.Context, thingId string, operatorId string) error {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return err
	}
	bsonOperatorId, err := primitive.ObjectIDFromHex(operatorId)
	if err != nil {
		return err
	}
	_, err = service.ThingOperatorCollection(ctx).DeleteOne(ctx, bson.M{"thingId": bsonThingId, "operatorId": bsonOperatorId})
	return err
}

func (service *ThingOperatorService) IsAssociationUnique(ctx context.Context, newThingOperator *model.ThingOperator) bool {
	var thingOperator model.ThingOperator
	query := bson.M{"thingId": newThingOperator.ThingId, "operatorId": newThingOperator.OperatorId}
	err := service.ThingOperatorCollection(ctx).FindOne(ctx, query).Decode(&thingOperator)
	return err != nil
}

func (service *ThingOperatorService) ThingOperatorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
}

