package services

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
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

type UserServiceInterface interface {
	Create(context.Context, *model.User) (*mongo.InsertOneResult, error)
	FindByUserEmail(ctx context.Context, email string) (*model.User, error)
	FindByUserId(ctx context.Context, userId string) (*model.User, error)
	IsUserUnique(ctx context.Context, newUser *model.User) bool
	FindUsersByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.User, error)
	CreateToken(*gin.Context, *model.User) (string, error)
	HashPassword(password string) string
	CheckPasswordHash(password, hash string) bool	
}

type UserService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewUserService(db *mgo.Session, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	return service.UserCollection(ctx).InsertOne(ctx, user)
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := service.UserCollection(ctx).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (service *UserService) FindByUserId(ctx context.Context, userId string) (*model.User, error) {
	bsonUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = service.UserCollection(ctx).FindOne(ctx, bson.M{"_id": bsonUserId}).Decode(&user)
	return &user, err
}

func (service *UserService) IsUserUnique(ctx context.Context, newUser *model.User) bool {
	query := bson.M{"name": newUser.DisplayName, "email": newUser.Email, "organizationId": newUser.OrganizationId}
	err := service.UserCollection(ctx).FindOne(ctx, query)
	return err == nil
}

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.User, error) {
	var users []*model.User
	cursor, err := service.UserCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ============== Service Helper Method(s) ================

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.ID
	atClaims["organizationId"] = user.OrganizationId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	var expirationDate int = int(time.Now().Add(5 * time.Hour).Unix())
	c.SetCookie("Authorization", token, expirationDate, "/", "", false, true)
	return token, nil
}

func (service *UserService) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic("Hashing password failed")
	}
	return string(bytes)
}

func (service *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (service *UserService) UserCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("User")
}
