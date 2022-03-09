package organization

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
	FindByOrganizationId(context.Context, string) (*model.Organization, error)
	FindByOrganizationApiKey(context.Context, string) (*model.Organization, error)
}

type OrganizationService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewOrganizationService(db *mgo.Session, c *config.Configuration) OrganizationServiceInterface {
	return &OrganizationService{config: c, db: db}
}

func (service *OrganizationService) Create(ctx context.Context, organization *model.Organization) (*mongo.InsertOneResult, error) {
	organization.ID = primitive.NewObjectID()
	organization.ApiKey = uuid.NewString()
	return service.organizationCollection(ctx).InsertOne(ctx, organization)
}

func (service *OrganizationService) FindByOrganizationId(ctx context.Context, organizationId string) (*model.Organization, error) {
	bsonOrganizationId, err := primitive.ObjectIDFromHex(organizationId)
	if err != nil {
		return nil, err
	}

	var organization model.Organization
	err = service.organizationCollection(ctx).FindOne(ctx, bson.M{"_id": bsonOrganizationId}).Decode(&organization)
	return &organization, err
}

func (service *OrganizationService) FindByOrganizationApiKey(ctx context.Context, api_key string) (*model.Organization, error) {
	var organization model.Organization
	err := service.organizationCollection(ctx).FindOne(ctx, bson.M{"api_key": api_key}).Decode(&organization)
	return &organization, err
}

// ============== Service Helper Method(s) ================

func (service *OrganizationService) organizationCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("Organization")
}
