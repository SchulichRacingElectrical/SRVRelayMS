package controllers

import (
	"database-ms/databases"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type Organization struct {
	OrganizationId	*string	`json:"organizationId,omitempty"`
	Name						*string	`json:"name"`
	ApiKey					*string `json:"apiKey,omitempty"`
}

func PostOrganization(c *gin.Context) {
	// Parse the body
	var newOrganization Organization
	if err := c.BindJSON(&newOrganization); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// TODO: Generate an API key
	// Create the organization
	_, _, err := databases.Database.Client.
		Collection("organizations").
			Add(databases.Database.Context, newOrganization)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	
	c.Status(http.StatusOK)
}

func GetOrganizations(c *gin.Context) {
	snapshot := databases.Database.Client.Collection("organizations")

	// Generate the organizaton array
	organizations := []Organization{}
	iter := snapshot.Documents(databases.Database.Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		organization := Organization{}
		organizationMap := doc.Data()
		delete(organizationMap, "apiKey")
		organizationMap["organizationId"] = doc.Ref.ID
		data, parseError := json.Marshal(organizationMap)
		if parseError != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		json.Unmarshal(data, &organization)
		organizations = append(organizations, organization)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	// TODO
}

func GetKeyVerification(c *gin.Context) {
	// TODO
}