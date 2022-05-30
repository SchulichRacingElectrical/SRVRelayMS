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
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Create the thing
		result := db.Create(&thing)
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of thing-operators
		thingOperators := []model.ThingOperator{}
		for _, operatorId := range thing.OperatorIds {
			thingOperator := model.ThingOperator{}
			thingOperator.ThingId = thing.Id
			thingOperator.OperatorId = operatorId
			thingOperators = append(thingOperators, thingOperator)
		}

		// Insert empty operatorIds
		if len(thing.OperatorIds) == 0 {
			thing.OperatorIds = []uuid.UUID{}
		}

		// Batch insert thing-operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Thing, *pgconn.PgError) {
	var things []*model.Thing
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the things associated with the given organization
		result := db.Where("organization_id = ?", organizationId).Find(&things)
		if result.Error != nil {
			return result.Error
		}

		// Get the ids of the relationship with operator
		for _, thing := range things {
			var thingOperators []*model.ThingOperator
			thing.OperatorIds = []uuid.UUID{}
			result = db.Table(model.TableNameThingOperator).Where("thing_id = ?", thing.Id).Find(&thingOperators)
			if result.Error != nil {
				return result.Error
			}
			for _, thingOperator := range thingOperators {
				thing.OperatorIds = append(thing.OperatorIds, thingOperator.OperatorId)
			}
		}
		return result.Error
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return things, nil
}

func (service *ThingService) Update(ctx context.Context, updatedThing *model.Thing) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Save the updated thing
		result := db.Updates(updatedThing)
		if result.Error != nil {
			return result.Error
		}

		// Delete all of the associated thing-operators
		result = db.Table(model.TableNameThingOperator).Where("thing_id = ?", updatedThing.Id).Delete(&model.ThingOperator{})
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of thing-operators
		thingOperators := []model.ThingOperator{}
		for _, operatorId := range updatedThing.OperatorIds {
			thingOperator := model.ThingOperator{}
			thingOperator.ThingId = updatedThing.Id
			thingOperator.OperatorId = operatorId
			thingOperators = append(thingOperators, thingOperator)
		}

		// Insert empty operatorIds
		if len(updatedThing.OperatorIds) == 0 {
			updatedThing.OperatorIds = []uuid.UUID{}
		}

		// Batch insert thing-operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
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
		return nil, utils.GetPostgresError(result.Error)
	}
	return thing, nil
}
