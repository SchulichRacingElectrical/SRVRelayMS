package controllers

import (
	"database-ms/databases"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, organization *map[string]interface{}) {
	c.Set("organization", *organization)
	c.Next()
}

// Returns the organization associated with the authorization
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handle API key
		apiKey := c.Request.Header.Get("apiKey")
		if apiKey != "" {
			organizations, err := databases.Firebase.Client.
				Collection("organizations").
					Where("apiKey", "==", apiKey).
						Documents(databases.Firebase.Context).
							GetAll()
			if err != nil || len(organizations) == 0 {
				c.Status(http.StatusUnauthorized)
				return
			}
			organization := organizations[0].Data()
			organization["organizationId"] = organizations[0].Ref.ID
			success(c, &organization)
			return
		} 

		// Handle regular authorization // TODO: Send permissions too?
		reqToken := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			c.Status(http.StatusUnauthorized)
			return
		}
		reqToken = splitToken[1]
		token, err := databases.Firebase.Auth.VerifyIDToken(databases.Firebase.Context, reqToken)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}
		organizationId := token.Claims["organizationId"].(string)
		doc, err := databases.Firebase.Client.
			Collection("organizations").
				Doc(organizationId).
					Get(databases.Firebase.Context)
		organization := doc.Data()
		organization["organizationId"] = organizationId
		success(c, &organization)
	}
}