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

type ThingServiceInterface interface {
	// Public
	FindByOrganizationId(context.Context, uuid.UUID) ([]*model.Thing, *pgconn.PgError)
	Create(context.Context, *model.Thing) *pgconn.PgError
	Update(context.Context, *model.Thing) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(ctx context.Context, thingID uuid.UUID) (*model.Thing, *pgconn.PgError)
}

type ThingService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewThingService(db *gorm.DB, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

// PUBLIC FUNCTIONS

func (service *ThingService) Create(ctx context.Context, thing *model.Thing) *pgconn.PgError {
	result := service.db.Create(&thing)
	return utils.GetPostgresError(result.Error)
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Thing, *pgconn.PgError) {
	var things []*model.Thing
	result := service.db.Where("organization_id = ?", organizationId).Find(&things)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return things, nil
}

func (service *ThingService) Update(ctx context.Context, updatedThing *model.Thing) *pgconn.PgError {
	result := service.db.Updates(updatedThing)
	return utils.GetPostgresError(result.Error)
}

func (service *ThingService) Delete(ctx context.Context, thingId uuid.UUID) *pgconn.PgError {
	thing := model.Thing{Base: model.Base{Id: thingId}}
	result := service.db.Delete(&thing)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *ThingService) FindById(ctx context.Context, thingId uuid.UUID) (*model.Thing, *pgconn.PgError) {
	var thing *model.Thing
	result := service.db.Where("id = ?", thingId).First(&thing)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return thing, nil
}
