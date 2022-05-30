package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type OrganizationServiceInterface interface {
	// Public
	FindByOrganizationId(context.Context, uuid.UUID) (*model.Organization, *pgconn.PgError)
	FindAllOrganizations(context.Context) ([]*model.Organization, *pgconn.PgError)
	Create(context.Context, *model.Organization) (*mongo.InsertOneResult, *pgconn.PgError)
	UpdateKey(context.Context, *model.Organization) *pgconn.PgError
	Update(context.Context, *model.Organization) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindByOrganizationApiKey(context.Context, string) (*model.Organization, *pgconn.PgError)
}

type OrganizationService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewOrganizationService(db *gorm.DB, c *config.Configuration) OrganizationServiceInterface {
	return &OrganizationService{config: c, db: db}
}

// PUBLIC FUNCTIONS

func (service *OrganizationService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) (*model.Organization, *pgconn.PgError) {
	organization := model.Organization{}
	organization.Id = organizationId
	result := service.db.First(&organization)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return &organization, nil
}

func (service *OrganizationService) FindAllOrganizations(ctx context.Context) ([]*model.Organization, *pgconn.PgError) {
	var organizations []*model.Organization
	result := service.db.Find(&organizations)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return organizations, nil
}

func (service *OrganizationService) Create(ctx context.Context, organization *model.Organization) (*mongo.InsertOneResult, *pgconn.PgError) {
	organization.APIKey = uuid.NewString()
	result := service.db.Create(&organization)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return nil, nil
}

func (service *OrganizationService) UpdateKey(ctx context.Context, organization *model.Organization) *pgconn.PgError {
	organization.APIKey = uuid.NewString()
	result := service.db.Updates(&organization)
	if result.Error != nil {
		return utils.GetPostgresError(result.Error)
	}
	return nil
}

func (service *OrganizationService) Update(ctx context.Context, updatedOrganization *model.Organization) *pgconn.PgError {
	prev, perr := service.FindByOrganizationId(ctx, updatedOrganization.Id)
	if perr != nil {
		return perr
	}
	updatedOrganization.APIKey = prev.APIKey
	result := service.db.Updates(&updatedOrganization)
	if result.Error != nil {
		return utils.GetPostgresError(result.Error)
	}
	return nil
}

func (service *OrganizationService) Delete(ctx context.Context, organizationId uuid.UUID) *pgconn.PgError {
	organization := model.Organization{}
	organization.Id = organizationId
	result := service.db.Delete(&organization)
	if result.Error != nil {
		return utils.GetPostgresError(result.Error)
	}
	return nil
}

// PRIVATE FUNCTIONS

func (service *OrganizationService) FindByOrganizationApiKey(ctx context.Context, APIKey string) (*model.Organization, *pgconn.PgError) {
	organization := model.Organization{}
	organization.APIKey = APIKey
	result := service.db.First(&organization)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return &organization, nil
}
