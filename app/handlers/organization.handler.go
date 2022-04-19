package handlers

import (
	"database-ms/app/middleware"
	models "database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service services.OrganizationServiceInterface
}

func NewOrganizationAPI(organizationService services.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler {service: organizationService}
}

func (handler *OrganizationHandler) CreateOrganization(ctx *gin.Context) {
	var newOrganization models.Organization
	ctx.BindJSON(&newOrganization)
	result := make(map[string]interface{})

	// Ensure the org name is unique
	organizations, err := handler.service.FindAllOrganizations(ctx.Request.Context())
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, "")
	} else {
		for _, org := range *organizations {
			if org.Name == newOrganization.Name {
				utils.Response(ctx, http.StatusConflict, "Duplicate organization name.")
				return
			}
		}	
	}

	// Create the organization
	_, err = handler.service.Create(ctx.Request.Context(), &newOrganization)
	if err == nil {
		newOrganization.ApiKey = ""
		result = utils.SuccessPayload(newOrganization, "Successfully created organization.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.EntityCreationError)
		utils.Response(ctx, http.StatusBadRequest, result)
	}
}

func (handler *OrganizationHandler) GetOrganization(ctx *gin.Context) {
	organization, _ := ctx.Get("organization")
	result := utils.SuccessPayload(organization, "Successfully retrieved organization.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) GetOrganizations(ctx *gin.Context) {
	organizations, err := handler.service.FindAllOrganizations(ctx.Request.Context())
	if !middleware.IsSuperAdmin(ctx) {
		for _, organization := range *organizations {
			organization.ApiKey = ""
		}
	}
	if err == nil {
		result := utils.SuccessPayload(organizations, "Successfully retrieved organizations.")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result := utils.NewHTTPError(utils.OrganizationsNotFound)
		utils.Response(ctx, http.StatusNotFound, result)
	}
}

func (handler *OrganizationHandler) UpdateOrganization(ctx *gin.Context) {
	// TODO: Only the admin of the org can update
}

func (handler *OrganizationHandler) DeleteOrganization(ctx *gin.Context) {
	// TODO: Only the super admin can do this
	// Will not do this for now
}
