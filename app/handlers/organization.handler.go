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
	if handler.service.IsOrganizationUnique(ctx, &newOrganization) {
		_, err := handler.service.Create(ctx.Request.Context(), &newOrganization)
		if err == nil {
			newOrganization.ApiKey = ""
			result := utils.SuccessPayload(newOrganization, "Successfully created organization.")
			utils.Response(ctx, http.StatusOK, result)
		} else {
			result := utils.NewHTTPError(utils.EntityCreationError)
			utils.Response(ctx, http.StatusBadRequest, result)
		}
	} else {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.OrganizationDuplicate))	
	}
}

func (handler *OrganizationHandler) GetOrganization(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization.ApiKey = ""
	}
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
	var updatedOrganization models.Organization
	ctx.BindJSON(&updatedOrganization)
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		if organization.ID == updatedOrganization.ID {
			err := handler.service.Update(ctx, &updatedOrganization)
			if err == nil {
				result := utils.SuccessPayload(nil, "Succesfully updated organization.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
	}
}

func (handler *OrganizationHandler) IssueNewAPIKey(ctx *gin.Context) {
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization, _ := middleware.GetOrganizationClaim(ctx)
		newKey, err := handler.service.UpdateKey(ctx, organization)
		if err == nil {
			result := utils.SuccessPayload(newKey, "Successfully created a new API key.")
			utils.Response(ctx, http.StatusOK, result)
		} else {
			utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
		}
	} else {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
	}
}

func (handler *OrganizationHandler) DeleteOrganization(ctx *gin.Context) {
	// TODO: Only the super admin can do this
	// Will not do this for now
}
