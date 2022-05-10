package services

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

type CommentServiceI interface {
	GetComments(context.Context, string, string) ([]*models.Comment, error)
	AddComment(context.Context, string, string, *models.Comment) error
	UpdateCommentContent(context.Context, string, *models.Comment) error
	DeleteComment(context.Context, string, string) error
}

type CommentService struct {
	config *config.Configuration
}

func NewCommentService(c *config.Configuration) CommentServiceI {
	return &CommentService{config: c}
}

func (service *CommentService) AddComment(ctx context.Context, entityType string, entityId string, comment *models.Comment) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonEntityId, err := primitive.ObjectIDFromHex(entityId)
	if err != nil {
		return err
	}

	// Check if entity exists
	res := database.Collection(utils.EntityCollectionTypes[entityType]).FindOne(ctx, bson.M{"_id": bsonEntityId})
	if res.Err() == mongo.ErrNoDocuments {
		return errors.New(utils.RunDNE)
	}

	comment.ID = primitive.NewObjectID()
	comment.CreationDate = utils.CurrentTimeInMilli()
	comment.Type = entityType
	comment.AssociatedId = bsonEntityId

	_, err = database.Collection("Comment").InsertOne(ctx, comment)
	return err
}

func (service *CommentService) GetComments(ctx context.Context, entityType string, entityId string) ([]*models.Comment, error) {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonEntityId, err := primitive.ObjectIDFromHex(entityId)
	if err != nil {
		return nil, err
	}

	// Get comments
	var comments []*models.Comment
	commentFilter := bson.M{"associatedId": bsonEntityId, "type": entityType}
	cursor, err := database.Collection("Comment").Find(ctx, commentFilter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (service *CommentService) UpdateCommentContent(ctx context.Context, commentId string, updatedComment *models.Comment) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}

	// Check if comment exists
	var comment models.Comment
	err = database.Collection("Comment").FindOne(ctx, bson.M{"_id": bsonCommentId}).Decode(&comment)
	if err != nil {
		return errors.New(utils.CommentDoesNotExist)
	}

	// Check if user owns the comment
	// TODO move this in handler
	if comment.UserID.Hex() != updatedComment.UserID.Hex() {
		return errors.New(utils.CommentCannotUpdateOtherUserComment)
	}

	_, err = database.Collection("Comment").UpdateOne(ctx,
		bson.M{"_id": bsonCommentId},
		bson.M{"$set": bson.M{
			"content":      updatedComment.Content,
			"creationDate": utils.CurrentTimeInMilli(),
		}})

	return err
}

func (service *CommentService) DeleteComment(ctx context.Context, commentId string, userId string) error {
	database, err := databases.GetDatabase(service.config.AtlasUri, service.config.MongoDbName, ctx)
	if err != nil {
		panic(err)
	}
	defer database.Client().Disconnect(ctx)

	bsonCommentId, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}

	// Check if comment exists
	var comment models.Comment
	err = database.Collection("Comment").FindOne(ctx, bson.M{"_id": bsonCommentId}).Decode(&comment)
	if err != nil {
		return errors.New(utils.CommentDoesNotExist)
	}

	// Check if user owns the comment
	// TODO Move to handler
	if comment.UserID.Hex() != userId {
		return errors.New(utils.CommentCannotUpdateOtherUserComment)
	}

	_, err = database.Collection("Comment").DeleteOne(ctx, bson.M{"_id": bsonCommentId})
	if err != nil {
		return err
	}

	return err
}
