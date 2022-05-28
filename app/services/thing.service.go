package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThingServiceInterface interface {
	Create(context.Context, *model.Thing) error
	FindByOrganizationId(context.Context, uuid.UUID) ([]*model.Thing, error)
	FindById(ctx context.Context, thingID uuid.UUID) (*model.Thing, error)
	Update(context.Context, *model.Thing) error
	Delete(context.Context, uuid.UUID) error
	AttachAssociatedOperatorIds(context.Context, *model.Thing)
}

type ThingService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewThingService(db *gorm.DB, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

func (service *ThingService) Create(ctx context.Context, thing *model.Thing) error {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Create the thing
		result := db.Create(&thing)
		if result.Error != nil {
			return result.Error
		}

		// Create ThingOperators for insertion
		thingOperators := []model.ThingOperator{}
		for _, operatorId := range thing.OperatorIds {
			thingOperator := model.ThingOperator{}
			thingOperator.ThingId = thing.Id
			thingOperator.OperatorId = operatorId
			thingOperators = append(thingOperators, thingOperator)
		}
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return utils.GetPostgresError(err)
	}
	return nil
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Thing, error) {
	var things []*model.Thing
	// cursor, err := service.ThingCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	// if err = cursor.All(ctx, &things); err != nil {
	// 	return nil, err
	// }
	// if things == nil {
	// 	things = []*model.Thing{}
	// } else {
	// 	for _, thing := range things {
	// 		service.AttachAssociatedOperatorIds(ctx, thing)
	// 	}
	// }
	return things, nil
}

func (service *ThingService) FindById(ctx context.Context, thingId uuid.UUID) (*model.Thing, error) {
	// var thing model.Thing
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }
	// err = service.ThingCollection(ctx).FindOne(ctx, bson.M{"_id": bsonThingId}).Decode(&thing)
	// if err == nil {
	// 	service.AttachAssociatedOperatorIds(ctx, &thing)
	// }
	// return &thing, err
	return nil, nil
}

func (service *ThingService) Update(ctx context.Context, updatedThing *model.Thing) error {
	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Create the thing
	// 	if _, err := db.Collection("Thing").UpdateOne(ctx, bson.M{"_id": updatedThing.ID}, bson.M{"$set": updatedThing}); err != nil {
	// 		return nil, err
	// 	}

	// 	// Find the current thing
	// 	currentThing, err := service.FindById(ctx, updatedThing.ID.Hex())
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Resolve which thing-operator relationships need to be inserted
	// 	var thingOperatorsToInsert []interface{}
	// 	for _, newOperatorId := range updatedThing.OperatorIds {
	// 		if !utils.IdInSlice(newOperatorId, currentThing.OperatorIds) {
	// 			thingOperatorsToInsert = append(thingOperatorsToInsert, bson.D{{"operatorId", newOperatorId}, {"thingId", updatedThing.ID}})
	// 		}
	// 	}
	// 	if len(thingOperatorsToInsert) > 0 {
	// 		if _, err = db.Collection("ThingOperator").InsertMany(ctx, thingOperatorsToInsert); err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	// Resolve which thing-operator relationships need to be deleted
	// 	var thingOperatorsToDelete []model.ThingOperator
	// 	for _, currentOperatorId := range currentThing.OperatorIds {
	// 		if !utils.IdInSlice(currentOperatorId, updatedThing.OperatorIds) {
	// 			thingOperatorsToDelete = append(thingOperatorsToDelete, model.ThingOperator{OperatorId: currentOperatorId, ThingId: updatedThing.ID})
	// 		}
	// 	}
	// 	for _, thingOperator := range thingOperatorsToDelete {
	// 		if _, err = db.Collection("ThingOperator").DeleteOne(ctx, bson.M{"operatorId": thingOperator.OperatorId, "thingId": thingOperator.ThingId}); err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	return nil, nil
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

func (service *ThingService) Delete(ctx context.Context, thingId uuid.UUID) error {
	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return err
	// }

	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)
	// 	if _, err := db.Collection("ThingOperator").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
	// 		return nil, err
	// 	}
	// 	if _, err := db.Collection("Sensor").DeleteMany(ctx, bson.M{"thingId": bsonThingId}); err != nil {
	// 		return nil, err
	// 	}
	// 	if _, err := db.Collection("Thing").DeleteOne(ctx, bson.M{"_id": bsonThingId}); err != nil {
	// 		return nil, err
	// 	}
	// 	// TODO: There will be a lot more things to delete in the future...
	// 	// Will need to delete associated runs and sessions
	// 	// Will need to delete associated presets
	// 	return nil, nil
	// }

	// if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
	// 	return err
	// }

	return nil
}

func (service *ThingService) AttachAssociatedOperatorIds(ctx context.Context, thing *model.Thing) {
	// thing.OperatorIds = []primitive.ObjectID{}
	// var thingOperators []*model.ThingOperator
	// dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return
	// }
	// thingOperatorCollection := dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
	// cursor, err := thingOperatorCollection.Find(ctx, bson.M{"thingId": thing.ID})
	// if err = cursor.All(ctx, &thingOperators); err != nil {
	// 	return
	// }
	// var operatorIds []primitive.ObjectID
	// for _, thingOperator := range thingOperators {
	// 	operatorIds = append(operatorIds, thingOperator.OperatorId)
	// }
	// if len(operatorIds) == 0 {
	// 	thing.OperatorIds = []primitive.ObjectID{}
	// } else {
	// 	thing.OperatorIds = operatorIds
	// }
}
