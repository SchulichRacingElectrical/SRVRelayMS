package run

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RunServiceI interface {
	Create(context.Context, *models.Run) error
	GetRun(context.Context, string) (*models.Run, error)
	GetComments(context.Context, string) ([]*models.Comment, error)
	AddComment(context.Context, string, *models.Comment) error
	UpdateContent(context.Context, string, *models.Comment) error
	DeleteComment(context.Context, string, string) error
	Delete(context.Context, string) error
}

type RunService struct {
	config *config.Configuration
}

func NewRunService(c *config.Configuration) RunServiceI {
	return &RunService{config: c}
}

func (service *RunService) Create(ctx context.Context, run *models.Run) error {
	run.ID = primitive.NewObjectID()
	_, err := service.getCollection(ctx, "Run").InsertOne(ctx, run)
	return err
}

func (service *RunService) GetRun(ctx context.Context, runId string) (*models.Run, error) {
	var run models.Run
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}
	err = service.getCollection(ctx, "Run").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&run)
	return &run, err
}

func (service *RunService) AddComment(ctx context.Context, runId string, comment *models.Comment) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return err
	}

	comment.ID = primitive.NewObjectID()
	comment.CreationDate = utils.CurrentTimeInMilli()

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Add to run comment id list
		insertUpdate := bson.M{"$push": bson.M{"commentIds": comment.ID}}
		if _, err := db.Collection("Run").UpdateByID(ctx, bsonRunId, insertUpdate); err != nil {
			return nil, err
		}

		// Create comment document
		_, err := db.Collection("Comment").InsertOne(ctx, comment)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := databases.WithTransaction(client, ctx, callback); err != nil {
		return err
	}

	return nil
}

func (service *RunService) GetComments(ctx context.Context, runId string) ([]*models.Comment, error) {
	// Get run
	var run models.Run
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	if err = service.getCollection(ctx, "Run").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&run); err != nil {
		return nil, err
	}

	// Get comments
	var comments []*models.Comment
	commentFilter := bson.M{"_id": bson.M{"$in": run.CommentsId}}
	cursor, err := service.getCollection(ctx, "Comment").Find(ctx, commentFilter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (service *RunService) UpdateContent(ctx context.Context, commentId string, comment *models.Comment) error {
	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}
	comment.CreationDate = utils.CurrentTimeInMilli()
	_, err = service.getCollection(ctx, "Comment").UpdateOne(ctx, bson.M{"_id": bsonCommentId}, bson.M{"$set": comment})
	return err
}

func (service *RunService) DeleteComment(ctx context.Context, runId string, commentId string) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return err
	}

	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Remove the comment id from Run comments Id list
		runUpdate := bson.M{"$pull": bson.M{"commentIds": bsonCommentId}}
		if _, err := db.Collection("Run").UpdateByID(ctx, bsonRunId, runUpdate); err != nil {
			return nil, err
		}

		// Delete comment
		_, err := db.Collection("Comment").DeleteOne(ctx, bson.M{"_id": bsonCommentId})
		return nil, err
	}

	_, err = databases.WithTransaction(client, ctx, callback)

	return err
}

func (service *RunService) Delete(ctx context.Context, runId string) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return err
	}

	client, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		return err
	}

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := client.Database(service.config.MongoDbName)

		// Get related comments
		runProjection := bson.D{{"commentIds", 1}}
		runFilter := bson.M{"_id": bsonRunId}
		runOpts := options.FindOne().SetProjection(runProjection)
		var runEntity map[string]interface{}
		if err := db.Collection("Run").FindOne(ctx, runFilter, runOpts).Decode(&runEntity); err != nil {
			return nil, err
		}

		// Delete related comments
		commentFilter := bson.M{"_id": bson.M{"$in": runEntity["commentIds"]}}
		if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
			return nil, err
		}

		// Delete run
		_, err := db.Collection("Run").DeleteOne(ctx, runFilter)
		return nil, err
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *RunService) getCollection(ctx context.Context, collection string) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}

	return dbClient.Database(service.config.MongoDbName).Collection(collection)
}
