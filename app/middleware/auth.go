package middleware

import (
	"context"
	"database-ms/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	organizationSrv "database-ms/app/services/organization"
	userSrv "database-ms/app/services/user"

	"gopkg.in/mgo.v2"
)

// Returns the organization associated with the authorization
func AuthorizationMiddleware(conf *config.Configuration, dbSession *mgo.Session) gin.HandlerFunc {

	userInterface := userSrv.NewUserService(dbSession, conf)
	organizationInterface := organizationSrv.NewOrganizationService(dbSession, conf)

	return func(c *gin.Context) {

		// Check Admin API Key
		apiToken := c.Request.Header.Get("api_token")
		if apiToken == conf.AdminKey {
			c.Set("admin", true)
			c.Next()
			return
		}

		// Check JWT token
		tokenString, err := c.Cookie("Authorization")
		if tokenString == "" || err != nil {
			respondWithError(c, http.StatusUnauthorized, "No authorization detected.")
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
			user, err := userInterface.FindByUserId(context.TODO(), userId)
			if err != nil {
				respondWithError(c, http.StatusUnauthorized, "User not found.")
				return
			}
			organization, err := organizationInterface.FindByOrganizationId(context.TODO(), user.OrganizationId)
			if err != nil {
				respondWithError(c, http.StatusUnauthorized, "Organization not found.")
				return
			}
			c.Set("admin", false)
			c.Set("user", user)
			c.Set("organization", organization)
		} else {
			fmt.Println(err)
			respondWithError(c, http.StatusInternalServerError, "Decrypting Token Failed.")
		}
	}
}

func respondWithError(c *gin.Context, httpErrorCode int, message string) {
	c.Status(httpErrorCode)
	c.AbortWithStatusJSON(httpErrorCode, message)
}
