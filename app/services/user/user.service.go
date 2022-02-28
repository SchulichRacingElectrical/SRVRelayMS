package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserServiceInterface interface {
	Create(context.Context, *model.User) error
	FindByUserEmail(context.Context, string) (*model.User, error)
	FindByUserId(context.Context, string) (*model.User, error)
	Delete(context.Context, string) error

	CreateToken(*gin.Context, *model.User) (string, error)
}

type UserService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewUserService(db *mgo.Session, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) error {
	return service.addUser(ctx, user)
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	return service.getUser(ctx, bson.M{"email": email})
}

func (service *UserService) FindByUserId(ctx context.Context, userId string) (*model.User, error) {

	return service.getUser(ctx, bson.M{"_id": bson.ObjectIdHex(userId)})
}

func (service *UserService) Update(ctx context.Context, userId string, user *model.User) error {
	query := bson.M{"_id": bson.ObjectIdHex(userId)}
	CustomBson := &utils.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return err
	}

	return service.collection().Update(query, change)
}

func (service *UserService) Delete(ctx context.Context, userId string) error {

	// TODO Delete userId from Thing UserId list

	return service.collection().RemoveId(bson.ObjectIdHex(userId))
}

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.ID
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	c.SetCookie("Authorization", token, -1, "/", "", false, true)
	return token, nil
}

// ============== Service Helper Method(s) ================

// ============== Common DB Operations ===================

func (service *UserService) addUser(ctx context.Context, user *model.User) error {
	return service.collection().Insert(user)
}

func (service *UserService) getUser(ctx context.Context, query interface{}) (*model.User, error) {
	var user model.User
	err := service.collection().Find(query).One(&user)
	return &user, err
}

func (service *UserService) getUsers(ctx context.Context, query interface{}) ([]*model.User, error) {
	var users []*model.User
	err := service.collection().Find(query).All(&users)
	return users, err
}

func (service *UserService) collection() *mgo.Collection {
	return service.db.DB(service.config.MongoDbName).C("User")
}
