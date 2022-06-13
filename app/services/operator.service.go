package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/app/utils"
	"database-ms/config"

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

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) *pgconn.PgError {
	result := service.db.Create(&operator)
	return utils.GetPostgresError(result.Error)
}

func (service *OperatorService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.Operator, *pgconn.PgError) {
	var operators = []*model.Operator{}
	result := service.db.Where("organization_id = ?", organizationId).Find(&operators)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return operators, nil
}

func (service *OperatorService) Update(ctx context.Context, updatedOperator *model.Operator) *pgconn.PgError {
	result := service.db.Updates(updatedOperator)
	return utils.GetPostgresError(result.Error)
}

func (service *OperatorService) Delete(ctx context.Context, operatorId uuid.UUID) *pgconn.PgError {
	operator := model.Operator{Base: model.Base{Id: operatorId}}
	result := service.db.Delete(&operator)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *OperatorService) FindById(ctx context.Context, operatorId uuid.UUID) (*model.Operator, *pgconn.PgError) {
	var operator *model.Operator
	result := service.db.Where("id = ?", operatorId).First(&operator)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return operator, nil
}
