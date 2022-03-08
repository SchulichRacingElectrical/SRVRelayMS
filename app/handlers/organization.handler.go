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
		res := &createEntityRes{
			ID: newOrganization.ID,
		}
		result = utils.SuccessPayload(res, "Successfully created organization")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *OrganizationHandler) Delete(c *gin.Context) {
	result := make(map[string]interface{})
	err := handler.organization.Delete(c.Request.Context(), c.Param("organizationId"))
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(c, http.StatusOK, result)
}
