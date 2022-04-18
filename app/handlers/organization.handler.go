package handlers

import (
	"database-ms/app/models"
	services "database-ms/app/services"
	"database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service services.OrganizationServiceInterface
}

func NewOrganizationAPI(organizationService services.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler {service: organizationService}
}

func (handler *OrganizationHandler) Create(ctx *gin.Context) {
	var newOrganization models.Organization
	ctx.BindJSON(&newOrganization)
	result := make(map[string]interface{})
	res, err := handler.service.Create(ctx.Request.Context(), &newOrganization)
	if err == nil {
		result = utils.SuccessPayload(res, "Successfully created organization.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(ctx, http.StatusBadRequest, result)
	}
}

func (handler *OrganizationHandler) GetOrganization(ctx *gin.Context) {
	result := make(map[string]interface{})
	organization, exists := ctx.Get("organization")
	if exists {
		result = utils.SuccessPayload(organization, "Successfully retrieved organization.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.OrganizationNotFound)
		utils.Response(ctx, http.StatusBadRequest, result)
	}
}

func (handler *OrganizationHandler) GetOrganizations(ctx *gin.Context) {
	result := make(map[string]interface{})
	organizations, err := handler.service.FindAllOrganizations(ctx.Request.Context())
	if err == nil {
		result = utils.SuccessPayload(organizations, "Successfully retrieved organizations.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.OrganizationsNotFound)
		utils.Response(ctx, http.StatusNotFound, result)
	}
}
