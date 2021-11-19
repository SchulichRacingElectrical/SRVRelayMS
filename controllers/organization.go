package controllers

import (
	"database-ms/databases"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type Organization struct {
	OrganizationId	*string	`json:"organizationId,omitempty"`
	Name						*string	`json:"name"`
	ApiKey					*string `json:"apiKey,omitempty"`
}

func GetOrganizations(c *gin.Context) {
	organizations := make([]interface{}, 0)
	iter := databases.Database.Client.
		Collection("organizations").
			Documents(databases.Database.Context)

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

func PostOrganization(c *gin.Context) {
	// Create the organization
	var newOrganization Organization
	if err := c.BindJSON(&newOrganization); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	key := uuid.New().String()
	newOrganization.ApiKey = &key
	var newOrganizationMap map[string]interface{}
	inrec, err := json.Marshal(newOrganization)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	json.Unmarshal(inrec, &newOrganizationMap)

	// Create firestore entry
	_, _, createError := databases.Database.Client.
		Collection("organizations").
			Add(databases.Database.Context, newOrganizationMap)
	if createError != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	
	c.Status(http.StatusOK)
}

func PutOrganization(c *gin.Context) {
	// TODO
}

func DeleteOrganization(c *gin.Context) {
	// TODO
}