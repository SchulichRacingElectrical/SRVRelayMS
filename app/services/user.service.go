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
	FindByUserEmail(context.Context, string) (*model.User, error)
	FindByUserId(context.Context, string) (*model.User, error)
	IsUserUnique(context.Context, *model.User) bool
	IsLastAdmin(context.Context, *model.User) (bool, error)
	FindUsersByOrganizationId(context.Context, primitive.ObjectID) ([]*model.User, error)
	Update(context.Context, *model.User) error
	Delete(context.Context, string) error
	CreateToken(*gin.Context, *model.User) (string, error)
	HashPassword(string) string
	CheckPasswordHash(string, string) bool	
}

type UserService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewUserService(db *mgo.Session, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	res, err := service.UserCollection(ctx).InsertOne(ctx, user)
	user.ID = (res.InsertedID).(primitive.ObjectID)
	return res, err
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
	users, err := service.FindUsersByOrganizationId(ctx, newUser.OrganizationId)
	if err == nil {
		for _, user := range users {
			// Email must be globally unique
			if newUser.Email == user.Email && newUser.ID != user.ID {
				return false
			}
			// Display name must be unique within the organization
			if newUser.DisplayName == user.DisplayName && newUser.OrganizationId == user.OrganizationId && newUser.ID != user.ID {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *UserService) IsLastAdmin(ctx context.Context, user *model.User) (bool, error) {
	users, err := service.FindUsersByOrganizationId(ctx, user.OrganizationId)
	if err == nil {
		for _, existingUser := range users {
			if user.ID != existingUser.ID && existingUser.Role == "Admin" {
				return false, nil
			}
		}
		return true, nil
	} else {
		return false, err
	}	
}

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) ([]*model.User, error) {
	var users []*model.User
	cursor, err := service.UserCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	for _, user := range users {
		user.Password = ""
	}
	if users == nil {
		users = []*model.User{}
	}
	return users, nil
}

func (service *UserService) Update(ctx context.Context, user *model.User) error {
	_, err := service.UserCollection(ctx).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (service *UserService) Delete(ctx context.Context, userId string) error {
	bsonUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	} else {
		_, err := service.UserCollection(ctx).DeleteOne(ctx, bson.M{"_id": bsonUserId})
		return err
	}
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
