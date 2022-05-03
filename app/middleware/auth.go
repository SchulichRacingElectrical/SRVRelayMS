package middleware

import (
	"context"
	"database-ms/config"
	"database-ms/utils"
	"fmt"
	"net/http"

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

	return func(ctx *gin.Context) {
		// Initialize admin flags to false.
		ctx.Set("super-admin", false)
		ctx.Set("org-admin", false)

		// Check API Key
		apiKey := ctx.Request.Header.Get("apiKey")

		// Check if API Key is the admin secret.
		switch apiKey {
			case "":
				break
			case conf.AdminKey:
				ctx.Set("super-admin", true)
				ctx.Next()
				return
			default:
				// Check if Api Key matches an organization.
				organization, err := organizationService.FindByOrganizationApiKey(context.TODO(), apiKey)
				if err != nil {
					utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
					return
				}
				
				// If an org is found, grant admin permissions on that org.
				ctx.Set("organization", organization)
				ctx.Set("org-admin", true)
				ctx.Next()
				return
		}

		// Check JWT token
		tokenString, err := ctx.Cookie("Authorization")
		if tokenString == "" || err != nil {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
			ctx.Abort()
			return
		}

		var hmacSampleSecret = []byte(conf.AccessSecret)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return hmacSampleSecret, nil
		})
		if err != nil {
			utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := fmt.Sprintf("%s", claims["userId"])
			user, err := userService.FindByUserId(context.TODO(), userId)
			if err != nil {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
				ctx.Abort()
				return
			}
			organization, err := organizationService.FindByOrganizationId(context.TODO(), user.OrganizationId)
			if err != nil {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
				ctx.Abort()
				return
			}
			ctx.Set("user", user)
			ctx.Set("organization", organization)
		} else {
			utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
			ctx.Abort()
		}
	}
}

func GetOrganizationClaim(ctx *gin.Context) (*models.Organization, error) {
	organizationInterface, organizationExists := ctx.Get("organization")
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
		return ctx.GetBool("org-admin")
	} else {
		return Roles[user.Role] >= Roles[role]	
	}
}

func IsSuperAdmin(ctx *gin.Context) bool {
	return ctx.GetBool("super-admin")
}