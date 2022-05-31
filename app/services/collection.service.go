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
