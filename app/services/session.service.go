package services

import (
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type SessionServiceI interface {
	CreateSession(context.Context, *models.Session) error
	FindById(context.Context, string) (*models.Session, error)
	GetSessionsByThingId(context.Context, string) ([]*models.Session, error)
	UpdateSession(context.Context, *models.Session) error
	DeleteSession(context.Context, string) error
}

type SessionService struct {
	config *config.Configuration
}

func NewSessionService(c *config.Configuration) SessionServiceI {
	return &SessionService{config: c}
}

func (service *SessionService) CreateSession(ctx context.Context, session *models.Session) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	session.ID = primitive.NewObjectID()
	// Check if Thing exists
	res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": session.ThingID})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("thing does not exist")
	}

	_, err = database.Collection("Session").InsertOne(ctx, session)
	return err
}

func (service *SessionService) FindById(ctx context.Context, sessionId string) (*models.Session, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	var session models.Session
	bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return nil, err
	}

	err = database.Collection("Session").FindOne(ctx, bson.M{"_id": bsonSessionId}).Decode(&session)
	if err == nil {
		return nil, err
	}

	return &session, err
}

func (service *SessionService) GetSessionsByThingId(ctx context.Context, thingId string) ([]*models.Session, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session
	cursor, err := database.Collection("Session").Find(ctx, bson.M{"thingId": bsonThingId})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (service *SessionService) UpdateSession(ctx context.Context, updatedSession *models.Session) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	// Check that start time < end time
	if updatedSession.StartDate > updatedSession.EndDate {
		return errors.New("startTime cannot be larger than Endtime")
	}

	// Check if Thing exists
	res := database.Collection("Thing").FindOne(ctx, bson.M{"_id": updatedSession.ThingID})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New("thing does not exist")
	}

	_, err = database.Collection("Session").UpdateOne(ctx, bson.M{"_id": updatedSession.ID}, bson.M{"$set": updatedSession})

	return err
}

func (service *SessionService) DeleteSession(ctx context.Context, sessionId string) error {
	bsonSessionId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Delete related comments
		commentFilter := bson.M{"associatedId": bsonSessionId, "type": utils.Session}
		if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
			return nil, err
		}

		// Set related run associatedId to empty
		runFilter := bson.D{{"sessionId", bsonSessionId}}
		runUpdate := bson.D{{"$set", bson.D{{"sessionId", nil}}}}
		if _, err := db.Collection("Run").UpdateMany(ctx, runFilter, runUpdate); err != nil {
			return nil, err
		}

		// Delete run
		_, err := db.Collection("Session").DeleteOne(ctx, bson.M{"_id": bsonSessionId})
		return nil, err
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}
