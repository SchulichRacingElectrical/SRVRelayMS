package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type CommentServiceInterface interface {
	// Public
	FindCommentsByContextId(context.Context, uuid.UUID) ([]*model.Comment, *pgconn.PgError)
	CreateComment(context.Context, *model.Comment) *pgconn.PgError
	UpdateComment(context.Context, *model.Comment) *pgconn.PgError
	DeleteComment(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.Comment, *pgconn.PgError)
}

type CommentService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewCommentService(db *gorm.DB, c *config.Configuration) CommentServiceInterface {
	return &CommentService{config: c, db: db}
}

// PUBLIC FUNCTIONS

func (service *CommentService) FindCommentsByContextId(ctx context.Context, contextId uuid.UUID) ([]*model.Comment, *pgconn.PgError) {
	var comments []*model.Comment
	result := service.db.Order("last_update desc").Find(&comments, "comment_id = ? AND (collection_id = ? OR session_id = ?)", nil, contextId, contextId)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return comments, nil
}

func (service *CommentService) CreateComment(ctx context.Context, comment *model.Comment) *pgconn.PgError {
	result := service.db.Create(&comment)
	return utils.GetPostgresError(result.Error)
}

func (service *CommentService) UpdateComment(ctx context.Context, comment *model.Comment) *pgconn.PgError {
	result := service.db.Updates(&comment)
	return utils.GetPostgresError(result.Error)
}

func (service *CommentService) DeleteComment(ctx context.Context, commentId uuid.UUID) *pgconn.PgError {
	comment := model.Comment{Base: model.Base{Id: commentId}}
	result := service.db.Delete(&comment)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *CommentService) FindById(ctx context.Context, commentId uuid.UUID) (*model.Comment, *pgconn.PgError) {
	var comment *model.Comment
	result := service.db.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return comment, nil
}
