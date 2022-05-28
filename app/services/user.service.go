package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	FindUsersByOrganizationId(context.Context, uuid.UUID) ([]*model.User, error)
	FindByUserEmail(context.Context, string) (*model.User, error)
	FindByUserId(context.Context, uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) (*mongo.InsertOneResult, error)
	IsUserUnique(context.Context, *model.User) bool
	IsLastAdmin(context.Context, *model.User) (bool, error)
	Update(context.Context, *model.User) error
	Delete(context.Context, uuid.UUID) error
	CreateToken(*gin.Context, *model.User) (string, error)
	HashPassword(string) string
	CheckPasswordHash(string, string) bool
}

type UserService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewUserService(db *gorm.DB, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	// res, err := service.UserCollection(ctx).InsertOne(ctx, user)
	// user.ID = (res.InsertedID).(primitive.ObjectID)
	// return res, err
	return nil, nil
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	// var user model.User
	// err := service.UserCollection(ctx).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	// return &user, err
	return nil, nil
}

func (service *UserService) FindByUserId(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	// bsonUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return nil, err
	// }
	// var user model.User
	// err = service.UserCollection(ctx).FindOne(ctx, bson.M{"_id": bsonUserId}).Decode(&user)
	// return &user, err
	return nil, nil
}

func (service *UserService) IsUserUnique(ctx context.Context, newUser *model.User) bool {
	// users, err := service.FindUsersByOrganizationId(ctx, newUser.OrganizationId)
	// if err == nil {
	// 	for _, user := range users {
	// 		// Email must be globally unique
	// 		if newUser.Email == user.Email && newUser.ID != user.ID {
	// 			return false
	// 		}
	// 		// Display name must be unique within the organization
	// 		if newUser.DisplayName == user.DisplayName && newUser.OrganizationId == user.OrganizationId && newUser.ID != user.ID {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// } else {
	// 	return false
	// }
	return false
}

func (service *UserService) IsLastAdmin(ctx context.Context, user *model.User) (bool, error) {
	// users, err := service.FindUsersByOrganizationId(ctx, user.OrganizationId)
	// if err == nil {
	// 	for _, existingUser := range users {
	// 		if user.ID != existingUser.ID && existingUser.Role == "Admin" {
	// 			return false, nil
	// 		}
	// 	}
	// 	return true, nil
	// } else {
	// 	return false, err
	// }
	return false, nil
}

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.User, error) {
	// var users []*model.User
	// cursor, err := service.UserCollection(ctx).Find(ctx, bson.D{{"organizationId", organizationId}})
	// if err = cursor.All(ctx, &users); err != nil {
	// 	return nil, err
	// }
	// for _, user := range users {
	// 	user.Password = ""
	// }
	// if users == nil {
	// 	users = []*model.User{}
	// }
	// return users, nil
	return nil, nil
}

func (service *UserService) Update(ctx context.Context, user *model.User) error {
	// _, err := service.UserCollection(ctx).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	// return err
	return nil
}

func (service *UserService) Delete(ctx context.Context, userId uuid.UUID) error {
	// bsonUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return err
	// } else {
	// 	_, err := service.UserCollection(ctx).DeleteOne(ctx, bson.M{"_id": bsonUserId})
	// 	return err
	// }
	return nil
}

// ============== Service Helper Method(s) ================

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.Id
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
