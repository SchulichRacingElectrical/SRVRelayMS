package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/app/utils"
	"database-ms/config"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

var tokenExpirationDuration time.Duration

type UserServiceInterface interface {
	// Public
	FindUsersByOrganizationId(context.Context, uuid.UUID) ([]*model.User, *pgconn.PgError)
	Create(context.Context, *model.User) *pgconn.PgError
	Update(context.Context, *model.User) *pgconn.PgError
	Delete(context.Context, uuid.UUID) *pgconn.PgError

	// Private
	FindByUserEmail(context.Context, string) (*model.User, *pgconn.PgError)
	FindByUserId(context.Context, uuid.UUID) (*model.User, *pgconn.PgError)
	IsLastAdmin(context.Context, *model.User) (bool, error)
	CreateToken(*gin.Context, *model.User) (string, error)
	BlacklistToken(*jwt.Token) error
	IsBlacklisted(*jwt.Token) bool
	HashPassword(string) string
	CheckPasswordHash(string, string) bool
}

type UserService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewUserService(db *gorm.DB, c *config.Configuration) UserServiceInterface {
	tokenExpirationDuration = 24 * time.Hour
	return &UserService{config: c, db: db}
}

// PUBLIC FUNCTIONS

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.User, *pgconn.PgError) {
	var users []*model.User
	result := service.db.Where("organization_id = ?", organizationId).Find(&users)
	if result.Error != nil {
		return nil, utils.GetPostgresError(result.Error)
	}
	return users, nil
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
	user := model.User{Base: model.Base{Id: userId}}
	result := service.db.Delete(&user)
	return utils.GetPostgresError(result.Error)
}

// PRIVATE FUNCTIONS

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, *pgconn.PgError) {
	var user *model.User
	result := service.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return user, nil
}

func (service *UserService) FindByUserId(ctx context.Context, userId uuid.UUID) (*model.User, *pgconn.PgError) {
	user := model.User{Base: model.Base{Id: userId}}
	result := service.db.First(&user)
	if result.Error != nil {
		return nil, &pgconn.PgError{}
	}
	return &user, nil
}

func (service *UserService) IsLastAdmin(ctx context.Context, user *model.User) (bool, error) {
	users, err := service.FindUsersByOrganizationId(ctx, user.OrganizationId)
	if err == nil {
		for _, existingUser := range users {
			if user.Id != existingUser.Id && existingUser.Role == "Admin" {
				return false, nil
			}
		}
		return true, nil
	} else {
		return false, err
	}
}

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	var expirationDate int = int(time.Now().Add(tokenExpirationDuration).Unix())
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.Id
	atClaims["organizationId"] = user.OrganizationId
	atClaims["exp"] = expirationDate
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	c.SetCookie("Authorization", token, expirationDate, "/", "", false, true)
	return token, nil
}

func (service *UserService) BlacklistToken(token *jwt.Token) error {
	// Extract expiration
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok != true {
		return errors.New("error parsing claims")
	}
	var exp time.Time
	exp = time.Unix(int64(claims["exp"].(float64)), 0)

	// Add token to blacklist table
	blacklist := model.Blacklist{
		Token:      token.Raw,
		Expiration: exp.Unix(),
	}
	result := service.db.Table(model.TableNameBlacklist).Create(&blacklist)
	return result.Error
}

func (service *UserService) IsBlacklisted(token *jwt.Token) bool {
	// Check if token exists in blacklist table
	count := int64(0)
	service.db.Table(model.TableNameBlacklist).Where("token = ?", token.Raw).Count(&count)
	return count > 0
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

func TokenExpirationDuration() int64 {
	return tokenExpirationDuration.Milliseconds()
}
