package handlers

import (
	"database-ms/app/middleware"
	model "database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service services.OrganizationServiceInterface
}

func NewOrganizationAPI(organizationService services.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{service: organizationService}
}

func (handler *OrganizationHandler) CreateOrganization(ctx *gin.Context) {
	// Attempt to extract the body
	var newOrganization model.Organization
	err := ctx.BindJSON(&newOrganization)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to create the organizaton
	_, err = handler.service.Create(ctx.Request.Context(), &newOrganization)
	if err != nil {
		result := utils.NewHTTPError(err.Error())
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}

	// Send the response
	newOrganization.APIKey = ""
	result := utils.SuccessPayload(newOrganization, "Successfully created organization.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) GetOrganization(ctx *gin.Context) {
	// Remove the API Key unless the requestor is an Admin
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		organization.APIKey = ""
	}

	// Send the response
	result := utils.SuccessPayload(organization, "Successfully retrieved organization.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) GetOrganizations(ctx *gin.Context) {
	// Attempt to read the organizations
	organizations, err := handler.service.FindAllOrganizations(ctx.Request.Context())
	if err != nil {
		result := utils.NewHTTPError(utils.OrganizationsNotFound)
		utils.Response(ctx, http.StatusNotFound, result)
		return
	}

	// Remove the API keys unless we are the super admin
	if !middleware.IsSuperAdmin(ctx) {
		for _, organization := range organizations {
			organization.APIKey = ""
		}
	}

	// Send the response
	result := utils.SuccessPayload(organizations, "Successfully retrieved organizations.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) UpdateOrganization(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedOrganization model.Organization
	err := ctx.BindJSON(&updatedOrganization)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin users
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against cross-tenant organization update
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if organization.Id != updatedOrganization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the organization
	err = handler.service.Update(ctx, &updatedOrganization)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Succesfully updated organization.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) IssueNewAPIKey(ctx *gin.Context) {
	// Guard against non-admin users
	if middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create a new key
	organization, _ := middleware.GetOrganizationClaim(ctx)
	err := handler.service.UpdateKey(ctx, organization)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPError(utils.InternalError))
		return
	}

	// Send the response
	result := utils.SuccessPayload(organization.APIKey, "Successfully created a new API key.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OrganizationHandler) DeleteOrganization(ctx *gin.Context) {
	// TODO: Only the super admin can do this
	// Will not do this for now
}
