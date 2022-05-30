package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorServiceInterface interface {
	Create(context.Context, *model.Operator) error
	FindById(context.Context, uuid.UUID) (*model.Operator, error)
	FindByOrganizationId(context.Context, uuid.UUID) ([]*model.Operator, error)
	Update(context.Context, *model.Operator) error
	Delete(context.Context, uuid.UUID) error
	IsOperatorUnique(context.Context, *model.Operator) bool
}

type OperatorService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewOperatorService(db *gorm.DB, c *config.Configuration) OperatorServiceInterface {
	return &OperatorService{db: db, config: c}
}

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) error {
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

		// Batch insert thing operators
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

func (service *OperatorService) FindById(ctx context.Context, operatorId uuid.UUID) (*model.Operator, error) {
	var operator *model.Operator
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the operator with the given id
		result := db.Where("id = ?", operatorId).First(&operator)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return operator, nil
}

func (service *OperatorService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Operator, error) {
	var operators = []*model.Operator{}
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Get the operators associated with the given organization
		result := db.Where("organization_id = ?", organizationId).Find(&operators)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, utils.GetPostgresError(err)
	}
	return operators, nil
}

func (service *OperatorService) Update(ctx context.Context, updatedOperator *model.Operator) error {
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

func (service *OperatorService) Delete(ctx context.Context, operatorId uuid.UUID) error {
	err := service.db.Transaction(func(db *gorm.DB) error {
		// Delete the specified operator
		operator := model.Operator{Base: model.Base{Id: operatorId}}
		result := db.Delete(&operator)
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

func (service *OperatorService) IsOperatorUnique(ctx context.Context, newOperator *model.Operator) bool {
	operators, err := service.FindByOrganizationId(ctx, newOperator.OrganizationId)
	if err == nil {
		for _, operator := range operators {
			if newOperator.Name == operator.Name && newOperator.Id != operator.Id {
				return false
			}
		}
		return true
	} else {
		return false
	}
}
