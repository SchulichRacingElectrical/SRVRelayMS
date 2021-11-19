package controllers

import (
	"database-ms/databases"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

// Returns the organization associated with the authorization
func Authorize(c *gin.Context) (map[string]interface{}, error){
	// Handle API key
	apiKey := c.Request.Header.Get("apiKey")
	if apiKey != "" {
		iter := databases.Database.Client.Collection("organizations").Documents(databases.Database.Context)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, errors.New("")
			}
			organization := doc.Data()
			if organization["apiKey"].(string) == apiKey {
				organization["organizationId"] = doc.Ref.ID
				return organization, nil
			}
		}
		return nil, errors.New("")
	} 

	// Handle regular authorization
	// TODO: Send permissions too?
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " Bearer")
	if len(splitToken) != 2 {
		return nil, errors.New("")
	}
	reqToken = splitToken[1]
	token, err := databases.Database.Auth.VerifyIDToken(databases.Database.Context, reqToken)
	if err != nil {
		return nil, errors.New("")
	}
	organizationId := token.Claims["organizationId"].(string)
	doc, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Database.Context)
	organization := doc.Data()
	organization["organizationId"] = organizationId
	return organization, nil
}