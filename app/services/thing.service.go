package services

import (
	"context"
	"database-ms/app/models"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"

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
	db     *mgo.Session
	config *config.Configuration
}

func NewThingService(db *mgo.Session, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

func (service *ThingService) Create(ctx context.Context, thing *model.Thing) error {
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		result, err := db.Collection("Thing").InsertOne(ctx, thing)
		if err != nil {
			return nil, err
		}
		thing.ID = (result.InsertedID).(primitive.ObjectID)
		var thingOperators []interface{}
		for _, operatorId := range thing.OperatorIds {
			thingOperators = append(thingOperators, bson.D{{"operatorId", operatorId}, {"thingId", thing.ID}})
		}
		if len(thingOperators) > 0 {
			if _, err = db.Collection("ThingOperator").InsertMany(ctx, thingOperators); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.Thing, error) {
	var things []*model.Thing
	cursor, err := service.ThingCollection(ctx).Find(ctx, bson.M{"organizationId": organizationId})
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
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Create the thing
		if _, err := db.Collection("Thing").UpdateOne(ctx, bson.M{"_id": updatedThing.ID}, bson.M{"$set": updatedThing}); err != nil {
			return nil, err
		}

		// Find the current thing
		currentThing, err := service.FindById(ctx, updatedThing.ID.Hex())
		if err != nil {
			return nil, err
		}

		// Resolve which thing-operator relationships need to be inserted
		var thingOperatorsToInsert []interface{}
		for _, newOperatorId := range updatedThing.OperatorIds {
			if !utils.IdInSlice(newOperatorId, currentThing.OperatorIds) {
				thingOperatorsToInsert = append(thingOperatorsToInsert, bson.D{{"operatorId", newOperatorId}, {"thingId", updatedThing.ID}})
			}
		}
		if len(thingOperatorsToInsert) > 0 {
			if _, err = db.Collection("ThingOperator").InsertMany(ctx, thingOperatorsToInsert); err != nil {
				return nil, err
			}
		}

		// Resolve which thing-operator relationships need to be deleted
		var thingOperatorsToDelete []model.ThingOperator
		for _, currentOperatorId := range currentThing.OperatorIds {
			if !utils.IdInSlice(currentOperatorId, updatedThing.OperatorIds) {
				thingOperatorsToDelete = append(thingOperatorsToDelete, model.ThingOperator{OperatorId: currentOperatorId, ThingId: updatedThing.ID})
			}
		}
		for _, thingOperator := range thingOperatorsToDelete {
			if _, err = db.Collection("ThingOperator").DeleteOne(ctx, bson.M{"operatorId": thingOperator.OperatorId, "thingId": thingOperator.ThingId}); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
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

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
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
		if _, err := db.Collection("RawDataPreset").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
			return nil, err
		}
		cursor, err := db.Collection("ChartPreset").Find(ctx, bson.M{"thingId": bsonThingId})
		if err == nil {
			var chartPresets []*models.ChartPreset
			if err = cursor.All(ctx, &chartPresets); err != nil {
				return nil, err
			} else {
				chartPresetIds := []primitive.ObjectID{}
				for _, preset := range chartPresets {
					chartPresetIds = append(chartPresetIds, preset.ID)
				}
				if _, err := db.Collection("Chart").DeleteMany(ctx, bson.M{"chartPresetId": bson.M{"$in": chartPresetIds}}); err != nil {
					return nil, err
				} else {
					if _, err := db.Collection("ChartPreset").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
						return nil, err
					}
				}
			}
		} else {
			return nil, err
		}

		// TODO: Delete associated runs and sessions and associated comments
		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *ThingService) IsThingUnique(ctx context.Context, newThing *model.Thing) bool {
	// TODO: Do with FindOne query rather than fetching everything
	things, err := service.FindByOrganizationId(ctx, newThing.OrganizationId)
	if err == nil {
		for _, thing := range things {
			if newThing.Name == thing.Name && newThing.ID != thing.ID {
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
	thingOperators := []*model.ThingOperator{}
	if cursor, err := service.ThingOperatorCollection(ctx).Find(ctx, bson.M{"thingId": thing.ID}); err == nil {
		if err = cursor.All(ctx, &thingOperators); err != nil {
			return
		}
		operatorIds := []primitive.ObjectID{}
		for _, thingOperator := range thingOperators {
			operatorIds = append(operatorIds, thingOperator.OperatorId)
		}
		if len(operatorIds) > 0 {
			thing.OperatorIds = operatorIds
		}
	} else {
		println(err.Error())
	}
}

func (service *ThingService) ThingCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Thing")
}

func (service *ThingService) ThingOperatorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
}
