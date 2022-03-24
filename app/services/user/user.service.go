package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

type UserServiceInterface interface {
	Create(context.Context, *model.User) (*mongo.InsertOneResult, error)
	FindByUserEmail(context.Context, string) (*model.User, error)
	FindByUserId(context.Context, string) (*model.User, error)

	CreateToken(*gin.Context, *model.User) (string, error)
}

type UserService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewUserService(db *mgo.Session, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	return service.userCollection(ctx).InsertOne(ctx, user)
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := service.userCollection(ctx).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (service *UserService) FindByUserId(ctx context.Context, userId string) (*model.User, error) {
	bsonUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = service.userCollection(ctx).FindOne(ctx, bson.M{"_id": bsonUserId}).Decode(&user)
	return &user, err
}

// User auth methods

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	var err error
	// Creating Access Token, do we need more claims?
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.ID
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	var expirationDate int = int(time.Now().Add(5 * time.Hour).Unix())
	c.SetCookie("Authorization", token, expirationDate, "/", "", false, true)
	return token, nil
}

// ============== Service Helper Method(s) ================

func (service *UserService) userCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("User")
}
