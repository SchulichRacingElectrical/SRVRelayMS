package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Organization struct {
	OrganizationId *string `json:"organizationId,omitempty"`
	Name           *string `json:"name" firestore:"name"`
	ApiKey         *string `json:"apiKey,omitempty" firestore:"apiKey"`
	// Default Thing
}

func GetOrganizations(c *gin.Context) {
	organizations := []Organization{}

	c.JSON(http.StatusOK, organizations)
}

func GetOrganization(c *gin.Context) {
	organization := Organization{}

	c.JSON(http.StatusOK, organization)
}

func PostOrganization(c *gin.Context) {

	c.Status(http.StatusOK)
}

func PutOrganization(c *gin.Context) {
	organization := Organization{}

	c.JSON(http.StatusOK, organization)
}

func DeleteOrganization(c *gin.Context) {

	c.Status(http.StatusOK)
}
