package services

import (
	"context"

	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

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
	return nil
}

func (service *ThingOperatorService) Delete(ctx context.Context, thingId string, operatorId string) error {
	return nil
}

func (service *ThingOperatorService) IsAssociationUnique(ctx context.Context, thingOperator *model.ThingOperator) bool {
	return false
}

func (service *ThingOperatorService) ThingOperatorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
}

