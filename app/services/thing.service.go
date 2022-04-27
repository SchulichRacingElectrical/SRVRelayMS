package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type ThingServiceInterface interface {
	Create(context.Context, *model.Thing) error
	FindByOrganizationId(context.Context, primitive.ObjectID) ([]*model.Thing, error)
	FindById(ctx context.Context, thingId string) (*model.Thing, error)
	Update(context.Context, *model.Thing) error
	Delete(context.Context, string) error
	IsThingUnique(context.Context, *model.Thing) bool
	AttachAssociatedOperatorIds(context.Context, *model.Thing)
}

type ThingService struct {
	db  	 *mgo.Session
	config *config.Configuration
}

func NewThingService(db *mgo.Session, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

func (service *ThingService) Create(ctx context.Context, thing *model.Thing) error {
	result, err := service.ThingCollection(ctx).InsertOne(ctx, thing)
	thing.ID = (result.InsertedID).(primitive.ObjectID)
	thing.OperatorIds = []primitive.ObjectID{}
	return err
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.Thing, error) {
	var things []*model.Thing
	cursor, err := service.ThingCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	if err = cursor.All(ctx, &things); err != nil {
		return nil, err
	}
	if things == nil {
		things = []*model.Thing{}
	} else {
		for _, thing := range things {
			service.AttachAssociatedOperatorIds(ctx, thing)
		}
	}
	return things, nil
}

func (service *ThingService) FindById(ctx context.Context, thingId string) (*model.Thing, error) {
	var thing model.Thing
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}
	err = service.ThingCollection(ctx).FindOne(ctx, bson.M{"_id": bsonThingId}).Decode(&thing)
	if err == nil {
		service.AttachAssociatedOperatorIds(ctx, &thing)
	}
	return &thing, err
}

func (service *ThingService) Update(ctx context.Context, updatedThing *model.Thing) error {
	if service.IsThingUnique(ctx, updatedThing) {
		_, err := service.ThingCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedThing.ID}, bson.M{"$set": updatedThing})
		return err
	} else {
		return errors.New("Thing name must remain unique.") // Could pass error code too?
	}
}

func (service *ThingService) Delete(ctx context.Context, thingId string) error {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		if _, err := db.Collection("ThingOperator").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
			return nil, err
		}
		if _, err := db.Collection("Sensor").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
			return nil, err
		} 
		if _, err := db.Collection("Thing").DeleteOne(ctx, bson.M{"_id": bsonThingId}); err != nil {
			return nil, err
		}
		// TODO: There will be a lot more things to delete in the future...
		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *ThingService) IsThingUnique(ctx context.Context, newThing *model.Thing) bool {
	things, err := service.FindByOrganizationId(ctx, newThing.OrganizationId)
	if err == nil {
		for _, thing := range things {
			if newThing.Name == thing.Name {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *ThingService) AttachAssociatedOperatorIds(ctx context.Context, thing *model.Thing) {
	thing.OperatorIds = []primitive.ObjectID{}
	var thingOperators []*model.ThingOperator
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return
	}
	thingOperatorCollection := dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
	cursor, err := thingOperatorCollection.Find(ctx, bson.M{"thingId": thing.ID})
	if err = cursor.All(ctx, &thingOperators); err != nil {
		return
	}	
	var operatorIds []primitive.ObjectID
	for _, thingOperator := range thingOperators {
		operatorIds = append(operatorIds, thingOperator.OperatorId)
	}
	if len(operatorIds) == 0 {
		thing.OperatorIds = []primitive.ObjectID{}
	} else {
		thing.OperatorIds = operatorIds
	}
}

func (service *ThingService) ThingCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Thing")
}
