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

type OperatorServiceInterface interface {
	// Public
	FindByOrganizationId(context.Context, uuid.UUID) ([]*model.Operator, *pgconn.PgError)
	Create(context.Context, *model.Operator) *pgconn.PgError
	Update(context.Context, *model.Operator) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindById(context.Context, uuid.UUID) (*model.Operator, *pgconn.PgError)
}

type OperatorService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewOperatorService(db *gorm.DB, c *config.Configuration) OperatorServiceInterface {
	return &OperatorService{db: db, config: c}
}

// PUBLIC FUNCTIONS

func (service *OperatorService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Operator, *pgconn.PgError) {
	var operators = []*model.Operator{}
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the operators associated with the given organization
		result := db.Where("organization_id = ?", organizationId).Find(&operators)
		if result.Error != nil {
			return result.Error
		}

		// Get the ids of the relationship with operator
		var thingOperators []*model.ThingOperator
		for _, operator := range operators {
			operator.ThingIds = []uuid.UUID{}
			result = db.Table(model.TableNameThingOperator).Where("operator_id = ?", operator.Id).Find(&thingOperators)
			if result.Error != nil {
				return result.Error
			}
			for _, thingOperator := range thingOperators {
				operator.ThingIds = append(operator.ThingIds, thingOperator.ThingId)
			}
		}
		return result.Error
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return operators, nil
}

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Create the operator
		result := db.Create(&operator)
		if result.Error != nil {
			return result.Error
		}

		// Generate the list of thing operators
		var thingOperators []model.ThingOperator
		for _, thingId := range operator.ThingIds {
			thingOperators = append(thingOperators, model.ThingOperator{
				OperatorId: operator.Id,
				ThingId:    thingId,
			})
		}

		// Insert empty thingIds
		if len(operator.ThingIds) == 0 {
			operator.ThingIds = []uuid.UUID{}
		}

		// Batch insert thing operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
}

func (service *OperatorService) Update(ctx context.Context, updatedOperator *model.Operator) *pgconn.PgError {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Save the updated operator
		db.Save(updatedOperator)

		// Delete all of the associated thing operators
		result := db.Table(model.TableNameThingOperator).Where("operator_id = ?", updatedOperator.Id).Delete(&model.ThingOperator{})
		if result.Error != nil {
			return result.Error
		}

		// Regenerate the list of thing-operators
		var thingOperators []model.ThingOperator
		for _, thingId := range updatedOperator.ThingIds {
			thingOperators = append(thingOperators, model.ThingOperator{
				OperatorId: updatedOperator.Id,
				ThingId:    thingId,
			})
		}

		// Insert empty thingIds
		if len(updatedOperator.ThingIds) == 0 {
			updatedOperator.ThingIds = []uuid.UUID{}
		}

		// Batch insert thing-operators
		result = db.Table(model.TableNameThingOperator).CreateInBatches(thingOperators, 100)
		return result.Error
	})
	return utils.GetPostgresError(err)
}

func (service *OperatorService) Delete(ctx context.Context, operatorId uuid.UUID) *pgconn.PgError {
	operator := model.Operator{Base: model.Base{Id: operatorId}}
	result := service.db.Delete(&operator)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *OperatorService) FindById(ctx context.Context, operatorId uuid.UUID) (*model.Operator, *pgconn.PgError) {
	var operator *model.Operator
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the operator with the given id
		result := db.Where("id = ?", operatorId).First(&operator)
		if result.Error != nil {
			return result.Error
		}
		// TODO: get the thingIds
		return result.Error
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return operator, nil
}
