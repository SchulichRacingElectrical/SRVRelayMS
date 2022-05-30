package services

import (
	"database-ms/app/model"
	"database-ms/config"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
)

type CollectionServiceInterface interface {
	// Public
	GetCollectionsByThingId(context.Context, uuid.UUID) ([]*model.Collection, *pgconn.PgError)
	CreateCollection(context.Context, *model.Collection) *pgconn.PgError
	UpdateCollection(context.Context, *model.Collection) *pgconn.PgError
	DeleteCollection(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.Collection, *pgconn.PgError)
}

type CollectionService struct {
	config *config.Configuration
}

func NewCollectionService(c *config.Configuration) CollectionServiceInterface {
	return &CollectionService{config: c}
}

// PUBLIC FUNCTIONS

func (service *CollectionService) GetCollectionsByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Collection, *pgconn.PgError) {

	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	// if err != nil {
	// 	return nil, err
	// }

	// var sessions []*models.Session
	// cursor, err := database.Collection("Session").Find(ctx, bson.M{"thingId": bsonThingId})
	// if err != nil {
	// 	return nil, err
	// }

	// if err = cursor.All(ctx, &sessions); err != nil {
	// 	return nil, err
	// }

	// return sessions, nil
	return nil, nil
}

func (service *CollectionService) CreateCollection(ctx context.Context, session *model.Collection) *pgconn.PgError {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// session.ID = primitive.NewObjectID()
	// // Check if Thing exists
	// res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": session.ThingID})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return primitive.NilObjectID, errors.New("thing does not exist")
	// }

	// result, err := database.Collection("Session").InsertOne(ctx, session)
	// if err != nil {
	// 	return primitive.NilObjectID, err
	// }
	// return result.InsertedID.(primitive.ObjectID), err
	return nil
}

func (service *CollectionService) UpdateCollection(ctx context.Context, updatedSession *model.Collection) *pgconn.PgError {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// // Check that start time < end time
	// if updatedSession.StartDate > updatedSession.EndDate {
	// 	return errors.New("startTime cannot be larger than Endtime")
	// }

	// // Check if Thing exists
	// res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": updatedSession.ThingID})
	// if res.Err() == mongo.ErrNoDocuments {
	// 	return errors.New("thing does not exist")
	// }

	// _, err = database.Collection("Session").UpdateOne(ctx, bson.M{"_id": updatedSession.ID}, bson.M{"$set": updatedSession})

	// return err
	return nil
}

func (service *CollectionService) DeleteCollection(ctx context.Context, sessionId uuid.UUID) *pgconn.PgError {
	// bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	// if err != nil {
	// 	return err
	// }

	// client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	// if err != nil {
	// 	return err
	// }
	// defer client.Disconnect(ctx)

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	db := client.Database(service.config.MongoDbName)

	// 	// Delete related comments
	// 	commentFilter := bson.M{"associatedId": bsonSessionId, "type": utils.Session}
	// 	if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
	// 		return nil, err
	// 	}

	// 	// Set related run associatedId to empty
	// 	runFilter := bson.D{{"sessionId", bsonSessionId}}
	// 	runUpdate := bson.D{{"$set", bson.D{{"sessionId", nil}}}}
	// 	if _, err := db.Collection("Run").UpdateMany(ctx, runFilter, runUpdate); err != nil {
	// 		return nil, err
	// 	}

	// 	// Delete run
	// 	_, err := db.Collection("Session").DeleteOne(ctx, bson.M{"_id": bsonSessionId})
	// 	return nil, err
	// }

	// _, err = databases.WithTransaction(client, ctx, callback)
	// return err
	return nil
}

// PRIVATE FUNCTIONS

func (service *CollectionService) FindById(ctx context.Context, sessionId uuid.UUID) (*model.Collection, *pgconn.PgError) {
	// database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer database.Client().Disconnect(ctx)

	// var session models.Session
	// bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	// if err != nil {
	// 	return nil, err
	// }

	// err = database.Collection("Session").FindOne(ctx, bson.M{"_id": bsonSessionId}).Decode(&session)
	// if err == nil {
	// 	return nil, err
	// }

	// return &session, err
	return nil, nil
}
