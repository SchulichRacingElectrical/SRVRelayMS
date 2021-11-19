package controllers

import (
	"database-ms/databases"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type Organization struct {
	OrganizationId	*string	`json:"organizationId,omitempty"`
	Name						*string	`json:"name"`
	ApiKey					*string `json:"apiKey,omitempty"`
}

func GetOrganizations(c *gin.Context) {
	snapshot := databases.Database.Client.Collection("organizations")

	// Generate the organizaton array
	organizations := make([]interface{}, 0)
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
		organization := doc.Data()
		delete(organization, "apiKey")
		organization["organizationId"] = doc.Ref.ID
		organizations = append(organizations, organization)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	organizationId := c.Param("organizationId")

	doc, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Database.Context)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	organization := doc.Data()
	organization["organizationId"] = doc.Ref.ID

	c.JSON(http.StatusOK, organization)
}

func GetKeyVerification(c *gin.Context) {
	// TODO
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

func PutOrganization(c *gin.Context) {
	// TODO
}