package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
)

type ThingServiceInterface interface {
	Create(context.Context, *model.Thing) error
	FindByOrganizationId(context.Context, string) ([]*model.Thing, error)
	FindById(context.Context, string) (*model.Thing, error)
	Update(context.Context, string, *model.ThingUpdate) error
	Delete(context.Context, string) error
}

type ThingService struct {
	db  	 *mgo.Session
	config *config.Configuration
}

func NewThingService(db *mgo.Session, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

func (service *ThingService) Create(ctx context.Context, thing *model.Thing) error {
	thing.ID = primitive.NewObjectID()
	_, err := service.thingCollection(ctx).InsertOne(ctx, thing)
	return err
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId string) ([]*model.Thing, error) {
	bsonOrganizationId, err := primitive.ObjectIDFromHex(organizationId)
	if err != nil {
		return nil, err
	}

	var things []*model.Thing
	cursor, err := service.thingCollection(ctx).Find(ctx, bson.D{{"organizationId", bsonOrganizationId}})
	if err = cursor.All(ctx, &things); err != nil {
		return nil, err
	}

	return things, nil
}

func (service *ThingService) FindById(ctx context.Context, thingId string) (*model.Thing, error) {
	var thing model.Thing
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	err = service.thingCollection(ctx).FindOne(ctx, bson.M{"_id": bsonThingId}).Decode(&thing)
	return &thing, err
}

func (service *ThingService) Update(ctx context.Context, thingId string, updates *model.ThingUpdate) error {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return err
	}
	_, err = service.thingCollection(ctx).UpdateOne(ctx, bson.M{"_id": bsonThingId}, bson.M{"$set": updates})
	return err
}

func (service *ThingService) Delete(ctx context.Context, thingId string) error {
	objectId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		thingCollection := db.Collection("Thing")

		// Get all related Sensor Id
		fmt.Println("Getting related sensors...")
		thingOpts := options.FindOne().SetProjection(bson.D{{"sensors", 1}})
		var thingEntity map[string]interface{}
		if err := thingCollection.FindOne(ctx, bson.M{"_id": objectId}, thingOpts).Decode(&thingEntity); err != nil {
			return nil, err
		}

		// Delete related sensor by using sensor id
		fmt.Println("Deleting related sensors...")
		sensorCollection := db.Collection("Sensor")
		sensorFilter := bson.M{"_id": bson.M{"$in": thingEntity["sensors"]}}
		// sensorOpts := options.Delete().SetHint(bson2.D{{"_id", 1}})
		if _, err := sensorCollection.DeleteMany(ctx, sensorFilter); err != nil {
			return nil, err
		}

		// Delete thing
		if _, err := thingCollection.DeleteOne(ctx, bson.M{"_id": objectId}); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *ThingService) thingCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Thing")
}
