package services

import (
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CollectionServiceInterface interface {
	CreateCollection(context.Context, *model.Collection) *pgconn.PgError
	FindById(context.Context, uuid.UUID) (*model.Collection, *pgconn.PgError)
	GetCollectionsByThingId(context.Context, uuid.UUID) ([]*model.Collection, *pgconn.PgError)
	UpdateCollection(context.Context, *model.Collection) *pgconn.PgError
	DeleteCollection(context.Context, uuid.UUID) *pgconn.PgError

	// Comments
	GetComments(context.Context, uuid.UUID) ([]*model.CollectionComment, error)
	GetComment(context.Context, uuid.UUID) (*model.CollectionComment, error)
	AddComment(context.Context, *model.CollectionComment) error
	UpdateCommentContent(context.Context, *model.CollectionComment) error
	DeleteComment(context.Context, uuid.UUID) error
}

type CollectionService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewCollectionService(db *gorm.DB, c *config.Configuration) CollectionServiceInterface {
	return &CollectionService{config: c, db: db}
}

func (service *CollectionService) CreateCollection(ctx context.Context, collection *model.Collection) *pgconn.PgError {
	result := service.db.Create(&collection)
	if result.Error != nil {
		var perr *pgconn.PgError
		errors.As(result.Error, &perr)
		return perr
	}
	return nil
}

func (service *CollectionService) FindById(ctx context.Context, collectionId uuid.UUID) (*model.Collection, *pgconn.PgError) {
	var collection *model.Collection
	result := service.db.Where("id = ?", collectionId).First(&collection)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return collection, nil
}

func (service *CollectionService) GetCollectionsByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Collection, *pgconn.PgError) {
	var collections []*model.Collection
	result := service.db.Where("thing_id = ?", thingId).Find(&collections)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return collections, nil
}

func (service *CollectionService) UpdateCollection(ctx context.Context, updatedCollection *model.Collection) *pgconn.PgError {
	var collection model.Collection
	if err := service.db.Where("id = ?", updatedCollection.Id).First(&collection).Error; err != nil {
		return &pgconn.PgError{}
	}

	result := service.db.Model(&collection).Updates(&updatedCollection)
	if result.Error != nil {
		return utils.GetPostgresError(result.Error)
	}
	return nil
}

func (service *CollectionService) DeleteCollection(ctx context.Context, collectionId uuid.UUID) *pgconn.PgError {
	collection := model.Collection{Base: model.Base{Id: collectionId}}
	result := service.db.Delete(&collection)
	return utils.GetPostgresError(result.Error)
}

func (service *CollectionService) GetComments(ctx context.Context, collectionId uuid.UUID) ([]*model.CollectionComment, error) {
	var comments []*model.CollectionComment
	result := service.db.Where("collection_id = ?", collectionId).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (service *CollectionService) GetComment(ctx context.Context, commentId uuid.UUID) (*model.CollectionComment, error) {
	var comment *model.CollectionComment
	result := service.db.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return comment, nil
}

func (service *CollectionService) AddComment(ctx context.Context, comment *model.CollectionComment) error {
	result := service.db.Create(&comment)
	return result.Error
}

func (service *CollectionService) UpdateCommentContent(ctx context.Context, updatedComment *model.CollectionComment) error {
	var comment model.CollectionComment
	if err := service.db.Where("id = ?", updatedComment.Id).First(&comment).Error; err != nil {
		return err
	}

	result := service.db.Model(&comment).Updates(&updatedComment)
	return result.Error
}

func (service *CollectionService) DeleteComment(ctx context.Context, commentId uuid.UUID) error {
	comment := model.CollectionComment{Base: model.Base{Id: commentId}}
	result := service.db.Delete(&comment)
	return result.Error
}
