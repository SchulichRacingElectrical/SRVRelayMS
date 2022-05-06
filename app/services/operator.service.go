package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type OperatorServiceInterface interface {
	Create(context.Context, *model.Operator) error
	FindById(context.Context, string) (*model.Operator, error)
	FindByOrganizationId(context.Context, primitive.ObjectID) ([]*model.Operator, error)
	Update(context.Context, *model.Operator) error
	Delete(context.Context, string) error
	IsOperatorUnique(context.Context, *model.Operator) bool
	AttachAssociatedThingIds(context.Context, *model.Operator)
}

type OperatorService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewOperatorService(db *mgo.Session, c *config.Configuration) OperatorServiceInterface {
	return &OperatorService{db: db, config: c}
}

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) error {
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		result, err := db.Collection("Operator").InsertOne(ctx, operator)
		if err != nil {
			return nil, err
		}
		operator.ID = (result.InsertedID).(primitive.ObjectID)
		var thingOperators []interface{}
		for _, thingId := range operator.ThingIds {
			thingOperators = append(thingOperators, bson.D{{"operatorId", operator.ID}, {"thingId", thingId}})
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

func (service *OperatorService) FindById(ctx context.Context, operatorId string) (*model.Operator, error) {
	var operator model.Operator
	bsonOperatorId, err := primitive.ObjectIDFromHex(operatorId)
	if err != nil {
		return nil, err
	}
	err = service.OperatorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonOperatorId}).Decode(&operator)
	if err == nil {
		service.AttachAssociatedThingIds(ctx, &operator)
	}
	return &operator, err
}

func (service *OperatorService) FindByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.Operator, error) {
	var operators []*model.Operator
	cursor, err := service.OperatorCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	if err = cursor.All(ctx, &operators); err != nil {
		return nil, err
	}
	if operators == nil {
		operators = []*model.Operator{}
	} else {
		for _, operator := range operators {
			service.AttachAssociatedThingIds(ctx, operator)
		}
	}
	return operators, nil
}

func (service *OperatorService) Update(ctx context.Context, updatedOperator *model.Operator) error {
	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Update the operator
		if _, err := db.Collection("Operator").UpdateOne(ctx, bson.M{"_id": updatedOperator.ID}, bson.M{"$set": updatedOperator}); err != nil {
			return nil, err
		}

		// Find the current operator
		currentOperator, err := service.FindById(ctx, updatedOperator.ID.Hex())
		if err != nil {
			return nil, err
		}

		// Resolve which thing-operator relationships to be inserted
		var thingOperatorsToInsert []interface{}
		for _, newThingId := range updatedOperator.ThingIds {
			if !utils.IdInSlice(newThingId, currentOperator.ThingIds) {
				thingOperatorsToInsert = append(thingOperatorsToInsert, bson.D{{"operatorId", updatedOperator.ID}, {"thingId", newThingId}})
			}
		}
		if len(thingOperatorsToInsert) > 0 {
			if _, err = db.Collection("ThingOperator").InsertMany(ctx, thingOperatorsToInsert); err != nil {
				return nil, err
			}
		}

		// Resolve which thing-operator relationships need to be deleted
		var thingOperatorsToDelete []model.ThingOperator
		for _, currentThingId := range currentOperator.ThingIds {
			if !utils.IdInSlice(currentThingId, updatedOperator.ThingIds) {
				thingOperatorsToDelete = append(thingOperatorsToDelete, model.ThingOperator{OperatorId: updatedOperator.ID, ThingId: currentThingId})
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

func (service *OperatorService) Delete(ctx context.Context, operatorId string) error {
	bsonOperatorId, err := primitive.ObjectIDFromHex(operatorId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		if _, err := db.Collection("ThingOperator").DeleteMany(ctx, bson.M{"operatorId": bsonOperatorId}); err != nil {
			return nil, err
		}
		if _, err := db.Collection("Operator").DeleteOne(ctx, bson.M{"_id": bsonOperatorId}); err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *OperatorService) IsOperatorUnique(ctx context.Context, newOperator *model.Operator) bool {
	operators, err := service.FindByOrganizationId(ctx, newOperator.OrganizationId)
	if err == nil {
		for _, operator := range operators {
			if newOperator.Name == operator.Name && newOperator.ID != operator.ID {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *OperatorService) AttachAssociatedThingIds(ctx context.Context, operator *model.Operator) {
	operator.ThingIds = []primitive.ObjectID{}
	var thingOperators []*model.ThingOperator
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return
	}
	thingOperatorCollection := dbClient.Database(service.config.MongoDbName).Collection("ThingOperator")
	cursor, err := thingOperatorCollection.Find(ctx, bson.M{"operatorId": operator.ID})
	if err = cursor.All(ctx, &thingOperators); err != nil {
		return
	}
	var thingIds []primitive.ObjectID
	for _, thingOperator := range thingOperators {
		thingIds = append(thingIds, thingOperator.ThingId)
	}
	if len(thingIds) == 0 {
		operator.ThingIds = []primitive.ObjectID{}
	} else {
		operator.ThingIds = thingIds
	}
}

func (service *OperatorService) OperatorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Operator")
}
