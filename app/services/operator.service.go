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
	db 			*mgo.Session
	config	*config.Configuration
}

func NewOperatorService(db *mgo.Session, c *config.Configuration) OperatorServiceInterface {
	return &OperatorService{db: db, config: c}
}

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) error {
	result, err := service.OperatorCollection(ctx).InsertOne(ctx, operator)
	operator.ID = (result.InsertedID).(primitive.ObjectID)
	operator.ThingIds = []primitive.ObjectID{}
	return err
}

func (service *OperatorService) FindById(ctx context.Context, operatorId string) (*model.Operator, error) {
	var operator model.Operator
	bsonOperatorId, err := primitive.ObjectIDFromHex(operatorId)
	if err != nil {
		return nil, err
	}
	err = service.OperatorCollection(ctx).FindOne(ctx, bson.M{"_id": bsonOperatorId}).Decode(&operator)
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
	if service.IsOperatorUnique(ctx, updatedOperator) {
		_, err := service.OperatorCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedOperator.ID}, bson.M{"$set": updatedOperator})
		return err
	} else {
		return errors.New("Operator name must remain unique.")
	}
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

	callback := func (sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)
		if _, err := db.Collection("ThingOperator").DeleteMany(ctx, bson.M{"operatorId": bsonOperatorId}); err != nil {
			return nil, err
		}
		if _, err := db.Collection("Operator").DeleteOne(ctx, bson.M{"_id": bsonOperatorId}); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *OperatorService) IsOperatorUnique(ctx context.Context, newOperator *model.Operator) bool {
	operators, err := service.FindByOrganizationId(ctx, newOperator.OrganizationId)
	if err == nil {
		for _, operator := range operators {
			if newOperator.Name == operator.Name {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *OperatorService) AttachAssociatedThingIds(ctx context.Context, operator *model.Operator) {
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
		thingIds = append(thingIds, thingOperator.ID)
	}
	operator.ThingIds = thingIds
}

func (service *OperatorService) OperatorCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Operator")
}