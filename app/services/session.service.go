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

	// Private
	FindById(context.Context, uuid.UUID) (*model.Session, *pgconn.PgError)
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

// PRIVATE FUNCTIONS

func (service *SessionService) FindById(ctx context.Context, sessionId uuid.UUID) (*model.Session, *pgconn.PgError) {
	var session *model.Session
	result := service.db.Where("id = ?", sessionId).First(&session)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return session, nil
}
