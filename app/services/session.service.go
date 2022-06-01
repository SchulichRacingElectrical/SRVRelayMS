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

type SessionServiceInterface interface {
	// Public - Session
	FindSessionsByThingId(context.Context, uuid.UUID) ([]*model.Session, *pgconn.PgError)
	CreateSession(context.Context, *model.Session) *pgconn.PgError
	UpdateSession(context.Context, *model.Session) *pgconn.PgError
	DeleteSession(context.Context, uuid.UUID) *pgconn.PgError

	// Public - Comments
	FindCommentsBySessionId(context.Context, uuid.UUID) ([]*model.SessionComment, *pgconn.PgError)
	CreateComment(context.Context, *model.SessionComment) *pgconn.PgError
	UpdateComment(context.Context, *model.SessionComment) *pgconn.PgError
	DeleteComment(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.Session, *pgconn.PgError)
	FindCommentById(context.Context, uuid.UUID) (*model.SessionComment, *pgconn.PgError)
}

type SessionService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewSessionService(db *gorm.DB, c *config.Configuration) SessionServiceInterface {
	return &SessionService{config: c, db: db}
}

// PUBLIC SESSION FUNCTIONS

func (service *SessionService) FindSessionsByThingId(ctx context.Context, thingId uuid.UUID) ([]*model.Session, *pgconn.PgError) {
	var sessions []*model.Session
	result := service.db.Where("thing_id = ?", thingId).Find(&sessions)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return sessions, nil
}

func (service *SessionService) CreateSession(ctx context.Context, session *model.Session) *pgconn.PgError {
	result := service.db.Create(&session)
	return utils.GetPostgresError(result.Error)
}

func (service *SessionService) UpdateSession(ctx context.Context, updatedSession *model.Session) *pgconn.PgError {
	result := service.db.Updates(&updatedSession)
	return utils.GetPostgresError(result.Error)
}

func (service *SessionService) DeleteSession(ctx context.Context, sessionId uuid.UUID) *pgconn.PgError {
	session := model.Session{Base: model.Base{Id: sessionId}}
	result := service.db.Delete(&session)
	return utils.GetPostgresError(result.Error)
}

// COMMENTS

func (service *SessionService) FindCommentsBySessionId(ctx context.Context, sessionId uuid.UUID) ([]*model.SessionComment, *pgconn.PgError) {
	var comments []*model.SessionComment
	result := service.db.Where("session_id = ?", sessionId).Find(&comments)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return comments, nil
}

func (service *SessionService) CreateComment(ctx context.Context, comment *model.SessionComment) *pgconn.PgError {
	result := service.db.Create(&comment)
	return utils.GetPostgresError(result.Error)
}

func (service *SessionService) UpdateComment(ctx context.Context, updatedComment *model.SessionComment) *pgconn.PgError {
	result := service.db.Updates(&updatedComment)
	return utils.GetPostgresError(result.Error)
}

func (service *SessionService) DeleteComment(ctx context.Context, commentId uuid.UUID) *pgconn.PgError {
	comment := model.SessionComment{Base: model.Base{Id: commentId}}
	result := service.db.Delete(&comment)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *SessionService) FindById(ctx context.Context, sessionId uuid.UUID) (*model.Session, *pgconn.PgError) {
	var session *model.Session
	result := service.db.Where("id = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return session, nil
}

func (service *SessionService) FindCommentById(ctx context.Context, commentId uuid.UUID) (*model.SessionComment, *pgconn.PgError) {
	var comment *model.SessionComment
	result := service.db.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return comment, nil
}
