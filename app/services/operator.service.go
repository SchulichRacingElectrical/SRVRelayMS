package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2"
)

type OperatorServiceInterface interface {
	Create(context.Context, *model.Operator) error
	FindByOrganizationId(context.Context, primitive.ObjectID) ([]*model.Operator, error)
	Update(context.Context, *model.Operator) error
	Delete(context.Context, string) error
}

type OperatorService struct {
	db 			*mgo.Session
	config	*config.Configuration
}

func NewOperatorService(db *mgo.Session, c *config.Configuration) OperatorServiceInterface {
	return &OperatorService{config: c, db: db}
}

func (service *OperatorService) Create(ctx context.Context, operator *model.Operator) error {
	// TODO
	return nil
}

func (service *OperatorService) FindByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.Operator, error) {
	// TODO
	return nil, nil
}

func (service *OperatorService) Update(ctx context.Context, updatedOperator *model.Operator) error {
	// TODO
	return nil
}

func (service *OperatorService) Delete(ctx context.Context, operatorId string) error {
	// TODO
	return nil
}