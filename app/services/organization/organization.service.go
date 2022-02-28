package organization

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"

	"github.com/google/uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OrganizationServiceInterface interface {
	Create(context.Context, *model.Organization) error
	FindByOrganizationId(context.Context, string) (*model.Organization, error)
	// FindUpdatedOrganization(context.Context, string, int64) ([]*model.Organization, error)
	// Update(context.Context, string, *model.OrganizationUpdate) error
	Delete(context.Context, string) error
}

type OrganizationService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewOrganizationService(db *mgo.Session, c *config.Configuration) OrganizationServiceInterface {
	return &OrganizationService{config: c, db: db}
}

func (service *OrganizationService) Create(ctx context.Context, organization *model.Organization) error {
	organization.ID = bson.NewObjectId()
	organization.ApiKey = uuid.NewString()
	return service.addOrganization(ctx, organization)
}

func (service *OrganizationService) FindByOrganizationId(ctx context.Context, organizationId string) (*model.Organization, error) {

	return service.getOrganization(ctx, bson.M{"thingId": bson.ObjectIdHex(organizationId)})

}

func (service *OrganizationService) Delete(ctx context.Context, organizationId string) error {

	// TODO Delete organizationId from Thing OrganizationId list

	return service.collection().RemoveId(bson.ObjectIdHex(organizationId))

}

// ============== Service Helper Method(s) ================

// ============== Common DB Operations ===================

func (service *OrganizationService) addOrganization(ctx context.Context, organization *model.Organization) error {
	return service.collection().Insert(organization)
}

func (service *OrganizationService) getOrganization(ctx context.Context, query interface{}) (*model.Organization, error) {
	var organization model.Organization
	err := service.collection().Find(query).One(&organization)
	return &organization, err
}

func (service *OrganizationService) getOrganizations(ctx context.Context, query interface{}) ([]*model.Organization, error) {
	var organizations []*model.Organization
	err := service.collection().Find(query).All(&organizations)
	return organizations, err
}

func (service *OrganizationService) collection() *mgo.Collection {
	return service.db.DB(service.config.MongoDbName).C("Organization")
}
