package handlers

import (
	"database-ms/app/models"
	organizationSrv "database-ms/app/services/organization"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	organization organizationSrv.OrganizationServiceInterface
}

func NewOrganizationAPI(organizationService organizationSrv.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{
		organization: organizationService,
	}
}

func (handler *OrganizationHandler) Create(c *gin.Context) {
	var newOrganization models.Organization
	c.BindJSON(&newOrganization)
	result := make(map[string]interface{})

	err := handler.organization.Create(c.Request.Context(), &newOrganization)
	var status int
	if err == nil {
		res := &createOrganizationRes{
			ID: newOrganization.ID,
		}
		result = utils.SuccessPayload(res, "Succesfully created organization")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *OrganizationHandler) GetOrganization(c *gin.Context) {
	// result := make(map[string]interface{})
	// organizations, err := handler.organization.GetOrganizations(c.Request.Context(), c.Param(""))
	// if err == nil {
	// 	result = utils.SuccessPayload(organizations, "Succesfully retrieved organizations")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	result = utils.NewHTTPError(utils.SensorsNotFound)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}

func (handler *OrganizationHandler) GetOrganizations(c *gin.Context) {
	result := make(map[string]interface{})
	organizations, err := handler.organization.GetOrganizations(c.Request.Context(), c.Param(""))
	if err == nil {
		result = utils.SuccessPayload(organizations, "Succesfully retrieved organizations")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}
