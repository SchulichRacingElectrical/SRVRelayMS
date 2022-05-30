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
	Create(context.Context, *model.Thing) *pgconn.PgError
	FindByOrganizationId(context.Context, uuid.UUID) ([]*model.Thing, error)
	FindById(ctx context.Context, thingID uuid.UUID) (*model.Thing, error)
	Update(context.Context, *model.Thing) error
	Delete(context.Context, uuid.UUID) error
}

type ThingService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewThingService(db *gorm.DB, c *config.Configuration) ThingServiceInterface {
	return &ThingService{config: c, db: db}
}

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

		// Batch insert thing-operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return utils.GetPostgresError(err)
	}
	return nil
}

func (service *ThingService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Thing, error) {
	var things []*model.Thing
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the things associated with the given organization
		result := db.Where("organization_id = ?", organizationId).Find(&things)
		if result.Error != nil {
			return result.Error
		}

		// Get the ids of the relationship with operator
		var thingOperators []*model.ThingOperator
		for _, thing := range things {
			result = db.Table(model.TableNameThingOperator).Where("thing_id = ?", thing).Find(&thingOperators)
			if result.Error != nil {
				return result.Error
			}
			for _, thingOperator := range thingOperators {
				thing.OperatorIds = append(thing.OperatorIds, thingOperator.OperatorId)
			}
		}
		return nil
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return things, nil
}

// Internal function
func (service *ThingService) FindById(ctx context.Context, thingId uuid.UUID) (*model.Thing, error) {
	var thing *model.Thing
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the thing with the given id
		result := db.Where("id = ?", thingId).First(&thing)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return thing, nil
}

func (service *ThingService) Update(ctx context.Context, updatedThing *model.Thing) error {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Save the updated thing
		db.Save(updatedThing)

		// Delete all of the associated thing-operators
		result := db.Table(model.TableNameThingOperator).Where("thing_id = ?", updatedThing.Id).Delete(&model.ThingOperator{})
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

		// Batch insert thing-operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return utils.GetPostgresError(err)
	}
	return nil
}

func (service *ThingService) Delete(ctx context.Context, thingId uuid.UUID) error {
	err := service.db.Transaction(func(db *gorm.DB) error { // Remove transaction
		// Delete the specified thing
		thing := model.Thing{Base: model.Base{Id: thingId}}
		result := db.Delete(&thing)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return utils.GetPostgresError(err)
	}
	return nil
}
