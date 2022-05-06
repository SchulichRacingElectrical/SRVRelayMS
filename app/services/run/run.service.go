package run

import (
	"context"
	"database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"database-ms/utils"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RunServiceI interface {
	CreateRun(context.Context, *models.Run) error
	FindById(context.Context, string) (*models.Run, error)
	GetRunsByThingId(context.Context, string) ([]*models.Run, error)
	UpdateRun(context.Context, *models.Run) error
	DeleteRun(context.Context, string) error

	GetComments(context.Context, string) ([]*models.Comment, error)
	AddComment(context.Context, string, *models.Comment) error
	UpdateCommentContent(context.Context, string, *models.Comment) error
	DeleteComment(context.Context, string) error
}

const RUN = "run"

type RunService struct {
	config *config.Configuration
}

func NewRunService(c *config.Configuration) RunServiceI {
	return &RunService{config: c}
}

func (service *RunService) CreateRun(ctx context.Context, run *models.Run) error {
	run.ID = primitive.NewObjectID()
	// Check if Thing exists
	res := service.getCollection(ctx, "Thing").FindOne(ctx, bson.M{"_id": run.ThingID})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return errors.New("thing does not exist")
	}

	// Check of Session exists
	res = service.getCollection(ctx, "Session").FindOne(ctx, bson.M{"_id": run.SessionId})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return errors.New("session does not exist")
	}

	_, err := service.getCollection(ctx, "Run").InsertOne(ctx, run)
	return err
}

func (services *RunService) FindById(ctx context.Context, runId string) (*models.Run, error) {
	var run models.Run
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	err = services.getCollection(ctx, "Run").FindOne(ctx, bson.M{"_id": bsonRunId}).Decode(&run)
	if err == nil {
		return nil, err
	}

	return &run, err
}

func (service *RunService) GetRunsByThingId(ctx context.Context, thingId string) ([]*models.Run, error) {
	bsonThingId, err := primitive.ObjectIDFromHex(thingId)
	if err != nil {
		return nil, err
	}

	var runs []*models.Run
	cursor, err := service.getCollection(ctx, "Run").Find(ctx, bson.M{"thingId": bsonThingId})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}

func (service *RunService) UpdateRun(ctx context.Context, updatedService *models.Run) error {
	// Check that start time < end time
	if updatedService.StartTime > updatedService.EndTime {
		return errors.New("startTime cannot be larger than Endtime")
	}

	// Check if Thing exists
	res := service.getCollection(ctx, "Thing").FindOne(ctx, bson.M{"_id": updatedService.ThingID})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return errors.New("thing does not exist")
	}

	// Check of Session exists
	res = service.getCollection(ctx, "Session").FindOne(ctx, bson.M{"_id": updatedService.SessionId})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return errors.New("session does not exist")
	}

	_, err := service.getCollection(ctx, "Comment").UpdateOne(ctx,
		bson.M{"_id": updatedService.ID},
		bson.M{"$set": updatedService})

	return err

}

func (service *RunService) DeleteRun(ctx context.Context, runId string) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
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
		commentFilter := bson.M{"associatedId": bsonRunId, "type": RUN}
		if _, err := db.Collection("Comment").DeleteMany(ctx, commentFilter); err != nil {
			return nil, err
		}

		// Delete run
		_, err := db.Collection("Run").DeleteOne(ctx, bson.M{"_id": bsonRunId})
		return nil, err
	}

	_, err = databases.WithTransaction(client, ctx, callback)
	return err
}

func (service *RunService) AddComment(ctx context.Context, runId string, comment *models.Comment) error {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return err
	}

	// Check if Run exists
	res := service.getCollection(ctx, "Run").FindOne(ctx, bson.M{"_id": bsonRunId})
	if res.Err().Error() == mongo.ErrNoDocuments.Error() {
		return errors.New("run does not exist")
	}

	comment.ID = primitive.NewObjectID()
	comment.CreationDate = utils.CurrentTimeInMilli()
	comment.Type = RUN
	comment.AssociatedId = bsonRunId

	_, err = service.getCollection(ctx, "Comment").InsertOne(ctx, comment)
	return err
}

func (service *RunService) GetComments(ctx context.Context, runId string) ([]*models.Comment, error) {
	bsonRunId, err := primitive.ObjectIDFromHex(runId)
	if err != nil {
		return nil, err
	}

	// Get comments
	var comments []*models.Comment
	commentFilter := bson.M{"associatedId": bsonRunId, "type": RUN}
	cursor, err := service.getCollection(ctx, "Comment").Find(ctx, commentFilter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (service *RunService) UpdateCommentContent(ctx context.Context, commentId string, updatedComment *models.Comment) error {
	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}

	// Check if comment exists
	var comment models.Comment
	err = service.getCollection(ctx, "Comment").FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err != nil {
		return errors.New(utils.CommentDoesNotExist)
	}

	// Check if user owns the comment
	if comment.UserID.Hex() != updatedComment.UserID.Hex() {
		return errors.New(utils.CommentCannotUpdateOtherUserComment)
	}

	_, err = service.getCollection(ctx, "Comment").UpdateOne(ctx,
		bson.M{"_id": bsonCommentId},
		bson.M{"$set": bson.M{
			"content":      updatedComment.Content,
			"creationDate": utils.CurrentTimeInMilli(),
		}})

	return err
}

func (service *RunService) DeleteComment(ctx context.Context, commentId string) error {
	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}

	res, err := service.getCollection(ctx, "Comment").DeleteOne(ctx, bson.M{"_id": bsonCommentId})
	if err != nil {
		return err
	}

	if res.DeletedCount == 1 {
		return errors.New(utils.CommentDoesNotExist)
	}

	return err
}

func (service *RunService) getCollection(ctx context.Context, collection string) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}

	return dbClient.Database(service.config.MongoDbName).Collection(collection)
}
