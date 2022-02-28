package middleware

import (
	"database-ms/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	model "database-ms/app/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Returns the organization associated with the authorization
func AuthorizationMiddleware(conf *config.Configuration, dbSession *mgo.Session) gin.HandlerFunc {

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
			user, err := getUserInfo(userId, conf, dbSession)
			if err != nil {
				respondWithError(c, http.StatusUnauthorized, "User not found.")
				return
			}
			organization, err := getOrganizationInfo(user.OrganizationId, conf, dbSession)
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

func getUserInfo(userId string, conf *config.Configuration, dbSession *mgo.Session) (*model.User, error) {
	var user model.User
	err := dbSession.DB(conf.MongoDbName).C("User").Find(bson.M{"_id": bson.ObjectIdHex(userId)}).One(&user)
	user.Password = ""
	return &user, err
}

func getOrganizationInfo(organizationId bson.ObjectId, conf *config.Configuration, dbSession *mgo.Session) (*model.Organization, error) {
	var organization model.Organization
	err := dbSession.DB(conf.MongoDbName).C("Organization").Find(bson.M{"_id": organizationId}).One(&organization)
	return &organization, err
}
