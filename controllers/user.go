package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserFirestore struct {
	StreamingSensors *[]string `json:"streamingSensors"`
}

type User struct {
	UserId         *string `json:"userId"`
	OrganizationId *string `json:"organizationId"`
	DisplayName    *string `json:"displayName"`
	Role           *string `json:"role"` // guest | admin | lead | member
	Email          *string `json:"email"`
	Disabled       *bool   `json:"disabled"`
}

func GetUsers(c *gin.Context) {

	users := []User{}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {

	response := User{}

	c.JSON(http.StatusOK, response)
}

func PostUser(c *gin.Context) {

	c.Status(http.StatusOK)
}

func PutUser(c *gin.Context) {

	c.Status(http.StatusOK)
}
