package controllers

import (
	"database-ms/databases"
	utils "database-ms/util"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type Organization struct {
	OrganizationId	*string	`json:"organizationId,omitempty"`
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
		delete(data, "apiKey")
		data["organizationId"] = doc.Ref.ID
		var organization Organization
		utils.JsonToStruct(data, &organization)
		organizations = append(organizations, organization)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)

	snapshot, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Database.Context)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	organization = snapshot.Data()
	organization["organizationId"] = organizationId
	c.JSON(http.StatusOK, organization)
}


type CreateOrganization struct {
	Name						*string	`json:"name" firestore:"name"`
	ApiKey					*string `json:"apiKey,omitempty" firestore:"apiKey"`
}

func PostOrganization(c *gin.Context) {
	var newOrganization CreateOrganization
	if err := c.BindJSON(&newOrganization); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	key := uuid.New().String()
	newOrganization.ApiKey = &key

	_, _, createError := databases.Database.Client.
		Collection("organizations").
			Add(databases.Database.Context, newOrganization)
	if createError != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	
	c.Status(http.StatusOK)
}


type PutOrganizationBody struct {
	Name						*string	`json:"name" binding:"required"`
	NewKey					*bool		`json:"newKey" binding:"required"`
}

func PutOrganization(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)

	updated := map[string]interface{}{}
	var request PutOrganizationBody
	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if request.NewKey != nil && *request.NewKey {
		key := uuid.New().String()
		updated["apiKey"] = key
	}
	if request.Name != nil {
		updated["name"] = *request.Name
	}
	
	_, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Set(databases.Database.Context, updated, firestore.MergeAll)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	newOrganization, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Database.Context)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, newOrganization.Data())
}

func DeleteOrganization(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)
	_, err := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Delete(databases.Database.Context)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}