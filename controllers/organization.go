package controllers

import (
	"database-ms/databases"
	utils "database-ms/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type Organization struct {
	OrganizationId	*string	`json:"organizationId,omitempty"`
	Name						*string	`json:"name" firestore:"name"`
	ApiKey					*string `json:"apiKey,omitempty" firestore:"apiKey"`
}

type PostOrganizationBody struct {
	Name						*string	`json:"name" firestore:"name"`
	ApiKey					*string `json:"apiKey,omitempty" firestore:"apiKey"`
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
		data := doc.Data()
		delete(data, "ApiKey")
		data["OrganizationId"] = doc.Ref.ID
		var organization Organization
		utils.JsonToStruct(data, &organization)
		organizations = append(organizations, organization)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	organization, exists := c.Get("organization")
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, organization)
}

func PostOrganization(c *gin.Context) {
	// Create the organization
	var newOrganization PostOrganizationBody
	if err := c.BindJSON(&newOrganization); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	key := uuid.New().String()
	newOrganization.ApiKey = &key

	// Create firestore entry
	_, _, createError := databases.Database.Client.
		Collection("organizations").
			Add(databases.Database.Context, newOrganization)
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