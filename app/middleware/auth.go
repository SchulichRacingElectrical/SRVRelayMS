package middleware

import (
	"context"
	"database-ms/config"
	"fmt"

	"database-ms/app/models"
	services "database-ms/app/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"gopkg.in/mgo.v2"
)

var Roles = map[string]int {
	"Admin": 4,
	"Lead": 3,
	"Member": 2,
	"Guest": 1,
	"Pending": 0,
}

// Returns the organization associated with the authorization
func AuthorizationMiddleware(conf *config.Configuration, dbSession *mgo.Session) gin.HandlerFunc {
	organizationService := services.NewOrganizationService(dbSession, conf)
	userService := services.NewUserService(dbSession, conf)

	return func(c *gin.Context) {
		// Initialize admin flags to false.
		c.Set("super-admin", false)
		c.Set("org-admin", false)

		// Check API Key
		apiKey := c.Request.Header.Get("apiKey")

		// Check if API Key is the admin secret.
		switch apiKey {
			case "":
				break
			case conf.AdminKey:
				c.Set("super-admin", true)
				c.Next()
				return
			default:
				// Check if Api Key matches an organization.
				organization, err := organizationService.FindByOrganizationApiKey(context.TODO(), apiKey)
				if err != nil {
					return
				}
				
				// If an org is found, grant admin permissions on that org.
				c.Set("organization", organization)
				c.Set("org-admin", true)
				c.Next()
				return
		}

		// Check JWT token
		tokenString, err := c.Cookie("Authorization")
		if tokenString == "" || err != nil {
			return
		}

		var hmacSampleSecret = []byte(conf.AccessSecret)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return hmacSampleSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := fmt.Sprintf("%s", claims["userId"])
			user, err := userService.FindByUserId(context.TODO(), userId)
			if err != nil {
				return
			}
			organization, err := organizationService.FindByOrganizationId(context.TODO(), user.OrganizationId)
			if err != nil {
				return
			}
			c.Set("user", user)
			c.Set("organization", organization)
		}
	}
}

func GetOrganizationClaim(ctx *gin.Context) (*models.Organization, error) {
	organizationInterface, organizationExists := ctx.Get("user")
	if organizationExists {
		return organizationInterface.(*models.Organization), nil
	} else {
		return nil, gin.Error{}
	} 
}

func GetUserClaim(ctx *gin.Context) (*models.User, error) {
	userInterface, userExists := ctx.Get("user")
	if userExists {
		return userInterface.(*models.User), nil
	} else {
		return nil, gin.Error{}
	}
}

func IsAuthorizationAtLeast(ctx *gin.Context, role string) bool {
	user, err := GetUserClaim(ctx)
	if err != nil {
		hasOrganizationKey, exists := ctx.Get("org-admin")
		return hasOrganizationKey.(bool) && exists
	} else {
		return Roles[user.Role] >= Roles[role]	
	}
}

func IsSuperAdmin(ctx *gin.Context) bool {
	superAdmin, exists := ctx.Get("super-admin")
	return superAdmin.(bool) && exists
}