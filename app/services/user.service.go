package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"database-ms/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	FindUsersByOrganizationId(context.Context, uuid.UUID) ([]*model.User, *pgconn.PgError)
	FindByUserEmail(context.Context, string) (*model.User, *pgconn.PgError)
	FindByUserId(context.Context, uuid.UUID) (*model.User, *pgconn.PgError)
	Create(context.Context, *model.User) *pgconn.PgError
	Update(context.Context, *model.User) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError
	IsLastAdmin(context.Context, *model.User) (bool, error)
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

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.User, *pgconn.PgError) {
	var users []*model.User
	result := service.db.Where("organization_id = ?", organizationId).Find(&users)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return users, nil
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, *pgconn.PgError) {
	var user *model.User
	result := service.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return user, nil
}

func (service *UserService) FindByUserId(ctx context.Context, userId uuid.UUID) (*model.User, *pgconn.PgError) {
	user := model.User{}
	user.Id = userId
	result := service.db.First(&user)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return &user, nil
}

func (service *UserService) Create(ctx context.Context, user *model.User) *pgconn.PgError {
	result := service.db.Create(&user)
	return utils.GetPostgresError(result.Error)
}

func (service *UserService) Update(ctx context.Context, user *model.User) *pgconn.PgError {
	result := service.db.Updates(&user)
	return utils.GetPostgresError(result.Error)
}

func (service *UserService) Delete(ctx context.Context, userId uuid.UUID) *pgconn.PgError {
	user := model.User{}
	user.Id = userId
	result := service.db.Delete(&user)
	return utils.GetPostgresError(result.Error)
}

func (service *UserService) IsLastAdmin(ctx context.Context, user *model.User) (bool, error) {
	users, err := service.FindUsersByOrganizationId(ctx, user.OrganizationId)
	if err == nil {
		for _, existingUser := range users {
			if user.Id != existingUser.Id && existingUser.Role == "Admin" {
				return false, nil
			}
		}
		print("here")
		return true, nil
	} else {
		return false, err
	}
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
