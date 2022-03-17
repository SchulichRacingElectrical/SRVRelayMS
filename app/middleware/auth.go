package middleware

import (
	"context"
	"database-ms/config"
	"fmt"

	organizationSrv "database-ms/app/services/organization"
	userSrv "database-ms/app/services/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"gopkg.in/mgo.v2"
)

// Returns the organization associated with the authorization
func AuthorizationMiddleware(conf *config.Configuration, dbSession *mgo.Session) gin.HandlerFunc {

	organizationService := organizationSrv.NewOrganizationService(dbSession, conf)
	userService := userSrv.NewUserService(dbSession, conf)

	return func(c *gin.Context) {

		// Initialize admin flags to false.
		c.Set("admin", false)
		c.Set("org-admin", false)

		// Check API Key
		apiKey := c.Request.Header.Get("apiKey")

		// Check if API Key is the admin secret.
		switch apiKey {
		case "":
			break
		case conf.AdminKey:
			c.Set("admin", true)
			c.Next()
			return
		default:
			// Check if Api Key matches an organization.
			organization, err := organizationService.FindByOrganizationApiKey(context.TODO(), apiKey)
			if err != nil {
				println(err.Error())
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
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := fmt.Sprintf("%s", claims["userId"])
			user, err := userService.FindByUserId(context.TODO(), userId)
			if err != nil {
				return
			}
			organization, err := organizationService.FindByOrganizationId(context.TODO(), user.OrganizationId.String())
			if err != nil {
				return
			}
			c.Set("user", user)
			c.Set("organization", organization)
		}
	}
}
