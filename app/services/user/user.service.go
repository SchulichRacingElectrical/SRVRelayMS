package user

import (
	"context"
	model "database-ms/app/models"
	"database-ms/config"
	"database-ms/utils"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserServiceInterface interface {
	Create(context.Context, *model.User) error
	FindByUserId(context.Context, string) (*model.User, error)
	Delete(context.Context, string) error
}

type UserService struct {
	db     *mgo.Session
	config *config.Configuration
}

func NewUserService(db *mgo.Session, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) Create(ctx context.Context, user *model.User) error {

	// find duplicate user

	// TODO add new user to Thing UserId list

	return service.addUser(ctx, user)
}

func (service *UserService) FindByUserId(ctx context.Context, userId string) (*model.User, error) {

	return service.getUser(ctx, bson.M{"thingId": bson.ObjectIdHex(userId)})

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
