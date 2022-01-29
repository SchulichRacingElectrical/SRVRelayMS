package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/utils"

	"github.com/google/uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OrganizationServiceInterface interface {
	Create(context.Context, *model.Organization) error
	// FindByThingId(context.Context, string) ([]*model.Organization, error)
	Update(context.Context, string, *model.Organization) error
	Delete(context.Context, string) error
	GetOrganizations(context.Context, string) ([]*model.Organization, error)
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

func (service *OrganizationService) FindUpdatedOrganization(ctx context.Context, thingId string, lastUpdate int64) ([]*model.Organization, error) {

	return service.getOrganizations(ctx, bson.M{
		"thingId": bson.ObjectIdHex(thingId),
		"lastUpdate": bson.M{
			"$gt": lastUpdate,
		},
	})

}

func (service *OrganizationService) Update(ctx context.Context, organizationId string, organization *model.Organization) error {
	query := bson.M{"_id": bson.ObjectIdHex(organizationId)}
	CustomBson := &utils.CustomBson{}
	change, err := CustomBson.Set(organization)
	if err != nil {
		return err
	}

	return service.collection().Update(query, change)
}

func (service *OrganizationService) Delete(ctx context.Context, organizationId string) error {

	// TODO Delete organizationId from Thing OrganizationId list
	return service.collection().RemoveId(bson.ObjectIdHex(organizationId))

}

func (service *OrganizationService) GetOrganizations(ctx context.Context, organizationId string) ([]*model.Organization, error) {
	return service.getOrganizations(ctx, bson.M{})
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
