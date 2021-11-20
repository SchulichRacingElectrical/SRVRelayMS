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
	// Default Thing
}

func GetOrganizations(c *gin.Context) {
	organizations := []Organization{}
	iter := databases.Firebase.Client.
		Collection("organizations").
			Documents(databases.Firebase.Context)

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
		organization := Organization{}
		utils.JsonToStruct(data, &organization)
		organizations = append(organizations, organization)
	}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)

	snapshot, err := databases.Firebase.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Firebase.Context)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	organization = snapshot.Data()
	organization["organizationId"] = organizationId
	c.JSON(http.StatusOK, organization)
}

type CreateOrganization struct {
	Name			*string	`json:"name" firestore:"name" binding:"required"`
	ApiKey		*string `json:"apiKey,omitempty" firestore:"apiKey"`
}

func PostOrganization(c *gin.Context) {
	var newOrganization CreateOrganization
	if err := c.BindJSON(&newOrganization); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	key := uuid.New().String()
	newOrganization.ApiKey = &key

	_, _, createError := databases.Firebase.Client.
		Collection("organizations").
			Add(databases.Firebase.Context, newOrganization)
	if createError != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	
	c.Status(http.StatusOK)
}

type PutOrganizationBody struct {
	Name			*string	`json:"name" binding:"required"`
	NewKey		*bool		`json:"newKey" binding:"required"`
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
	if *request.NewKey {
		key := uuid.New().String()
		updated["apiKey"] = key
	}
	updated["name"] = *request.Name
	
	_, err := databases.Firebase.Client.
		Collection("organizations").
			Doc(organizationId).
				Set(databases.Firebase.Context, updated, firestore.MergeAll)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	newOrganization, err := databases.Firebase.Client.
		Collection("organizations").
			Doc(organizationId).
				Get(databases.Firebase.Context)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, newOrganization.Data())
}

func DeleteOrganization(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)

	// Delete all of the associated users
	userIds := []string{}
	iter := databases.Firebase.Auth.Users(databases.Firebase.Context, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if user.CustomClaims["organizationId"] == organizationId {
			userIds = append(userIds, user.UID)
		}
	}
	_, deletionError := databases.Firebase.Auth.DeleteUsers(databases.Firebase.Context, userIds)
	if deletionError != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Delete the organization
	_, err := databases.Firebase.Client.
		Collection("organizations").
			Doc(organizationId).
				Delete(databases.Firebase.Context)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}