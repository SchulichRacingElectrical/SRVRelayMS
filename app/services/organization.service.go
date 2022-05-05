package services

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type OrganizationServiceInterface interface {
	Create(context.Context, *model.Organization) (*mongo.InsertOneResult, error)
	FindByOrganizationId(context.Context, primitive.ObjectID) (*model.Organization, error) // Should probably just cut this
	FindByOrganizationIdString(context.Context, string) (*model.Organization, error)
	FindByOrganizationApiKey(context.Context, string) (*model.Organization, error)
	FindAllOrganizations(context.Context) ([]*model.Organization, error)
	Update(context.Context, *model.Organization) error
	UpdateKey(context.Context, *model.Organization) (string, error)
	Delete(context.Context, primitive.ObjectID) error
	IsOrganizationUnique(context.Context, *model.Organization) bool
}

type OrganizationService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewOrganizationService(db *mgo.Session, c *config.Configuration) OrganizationServiceInterface {
	return &OrganizationService{config: c, db: db}
}

func (service *OrganizationService) Create(ctx context.Context, organization *model.Organization) (*mongo.InsertOneResult, error) {
	organization.ApiKey = uuid.NewString()
	res, err := service.OrganizationCollection(ctx).InsertOne(ctx, organization)
	organization.ID = (res.InsertedID).(primitive.ObjectID)
	return res, err
}

func (service *OrganizationService) FindByOrganizationIdString(ctx context.Context, organizationId string) (*model.Organization, error) {
	bsonOrganizationId, err := primitive.ObjectIDFromHex(organizationId)
	if err != nil {
		return nil, err
	}
	var organization model.Organization
	err = service.OrganizationCollection(ctx).FindOne(ctx, bson.M{"_id": bsonOrganizationId}).Decode(&organization)
	return &organization, err
}

func (service *OrganizationService) FindByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) (*model.Organization, error) {
	var organization model.Organization
	err := service.OrganizationCollection(ctx).FindOne(ctx, bson.M{"_id": organizationId}).Decode(&organization)
	return &organization, err
}

func (service *OrganizationService) FindByOrganizationApiKey(ctx context.Context, api_key string) (*model.Organization, error) {
	var organization model.Organization
	err := service.OrganizationCollection(ctx).FindOne(ctx, bson.M{"apiKey": api_key}).Decode(&organization)
	return &organization, err
}

func (service *OrganizationService) FindAllOrganizations(ctx context.Context) ([]*model.Organization, error) {
	var organizations []*model.Organization
	cursor, err := service.OrganizationCollection(ctx).Find(ctx, bson.D{})
	if err != nil || cursor.All(ctx, &organizations) != nil {
		return nil, err
	}
	if organizations == nil {
		organizations = []*model.Organization{}
	}
	return organizations, err
}

func (service *OrganizationService) Update(ctx context.Context, updatedOrganization *model.Organization) error {
	organization, err := service.FindByOrganizationIdString(ctx, updatedOrganization.ID.Hex())
	if err == nil {
		updatedOrganization.ApiKey = organization.ApiKey
		_, err = service.OrganizationCollection(ctx).UpdateOne(ctx, bson.M{"_id": updatedOrganization.ID}, bson.M{"$set": updatedOrganization})
		return err
	} else {
		return err
	}
}

func (service *OrganizationService) UpdateKey(ctx context.Context, organization *model.Organization) (string, error) {
	organization.ApiKey = uuid.NewString()	
	_, err := service.OrganizationCollection(ctx).UpdateOne(ctx, bson.M{"_id": organization.ID}, bson.M{"$set": organization})
	return organization.ApiKey, err
}

func (service *OrganizationService) Delete(ctx context.Context, organizationId primitive.ObjectID) error {
	// TODO
	// Transaction will be super nasty :(
	// Will need to delete all associated users, things, sensors, raw data presets, 
	// chart presets, runs, sessions, and comments!
	return nil
}

func (service *OrganizationService) IsOrganizationUnique(ctx context.Context, organization *model.Organization) bool {
	organizations, err := service.FindAllOrganizations(ctx)
	if err == nil {
		for _, org := range organizations {
			if organization.Name == org.Name {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// ============== Service Helper Method(s) ================

func (service *OrganizationService) OrganizationCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Organization")
}
