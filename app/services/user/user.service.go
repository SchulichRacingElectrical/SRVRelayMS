package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/databases"
	"fmt"
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
	FindUsersByOrganizationId(context.Context, primitive.ObjectID) (*[]model.User, error)

	CreateToken(*gin.Context, *model.User) (string, error)
	DecodeToken(context.Context, string) (*model.User, error)
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

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId primitive.ObjectID) (*[]model.User, error) {

	var users []model.User
	cursor, err := service.userCollection(ctx).Find(ctx, bson.M{"organizationId": organizationId})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	// Remove un from users list as it's a secret.
	for i := range users {
		users[i].Password = ""
	}

	return &users, err
}

// User auth methods

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	var err error
	// Creating Access Token, do we need more claims?
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.ID
	atClaims["name"] = user.DisplayName
	atClaims["email"] = user.Email
	atClaims["role"] = user.Roles

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	var expirationDate int = int(time.Now().Add(5 * time.Hour).Unix())
	c.SetCookie("Authorization", token, expirationDate, "/", "", false, true)
	return token, nil
}

func (service *UserService) DecodeToken(ctx context.Context, tokenString string) (*model.User, error) {
	var hmacSampleSecret = []byte(service.config.AccessSecret)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	var user model.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.DisplayName = fmt.Sprintf("%s", claims["name"])
		user.Email = fmt.Sprintf("%s", claims["email"])
		user.Roles = fmt.Sprintf("%s", claims["role"])
	} else {
		return nil, err
	}
	return &user, nil
}

// ============== Service Helper Method(s) ================

func (service *UserService) userCollection(ctx context.Context) *mongo.Collection {
	dbClient, err := databases.GetDBClient(service.config.AtlasUri, ctx)
	if err != nil {
		panic(err)
	}
	return dbClient.Database(service.config.MongoDbName).Collection("User")
}
