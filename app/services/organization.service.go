package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type OrganizationServiceInterface interface {
	FindByOrganizationId(context.Context, uuid.UUID) (*model.Organization, error)
	FindByOrganizationApiKey(context.Context, string) (*model.Organization, error)
	FindAllOrganizations(context.Context) ([]*model.Organization, error)
	Create(context.Context, *model.Organization) (*mongo.InsertOneResult, error)
	UpdateKey(context.Context, *model.Organization) error
	Update(context.Context, *model.Organization) error
	Delete(context.Context, uuid.UUID) error
}

type OrganizationService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewOrganizationService(db *gorm.DB, c *config.Configuration) OrganizationServiceInterface {
	return &OrganizationService{config: c, db: db}
}

func (service *OrganizationService) FindByOrganizationId(ctx context.Context, organizationId uuid.UUID) (*model.Organization, error) {
	organization := model.Organization{}
	organization.Id = organizationId
	result := service.db.First(&organization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &organization, nil
}

func (service *OrganizationService) FindByOrganizationApiKey(ctx context.Context, APIKey string) (*model.Organization, error) {
	organization := model.Organization{}
	organization.APIKey = APIKey
	result := service.db.First(&organization)
	if result.Error != nil {
		return nil, result.Error
	}
	return &organization, nil
}

func (service *OrganizationService) FindAllOrganizations(ctx context.Context) ([]*model.Organization, error) {
	var organizations []*model.Organization
	result := service.db.Find(&organizations)
	if result.Error != nil {
		return nil, result.Error
	}
	return organizations, nil
}

func (service *OrganizationService) Create(ctx context.Context, organization *model.Organization) (*mongo.InsertOneResult, error) {
	organization.APIKey = uuid.NewString()
	result := service.db.Create(&organization)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

func (service *OrganizationService) UpdateKey(ctx context.Context, organization *model.Organization) error {
	organization.APIKey = uuid.NewString()
	result := service.db.Save(&organization)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service *OrganizationService) Update(ctx context.Context, updatedOrganization *model.Organization) error {
	prev, err := service.FindByOrganizationId(ctx, updatedOrganization.Id)
	if err != nil {
		return err
	}
	updatedOrganization.APIKey = prev.APIKey
	result := service.db.Model(&updatedOrganization).Select("*").Updates(&updatedOrganization)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service *OrganizationService) Delete(ctx context.Context, organizationId uuid.UUID) error {
	organization := model.Organization{}
	organization.Id = organizationId
	result := service.db.Delete(&organization)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
