package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OperatorHandler struct {
	service services.OperatorServiceInterface
}

func NewOperatorAPI(operatorService services.OperatorServiceInterface) *OperatorHandler {
	return &OperatorHandler{service: operatorService}
}

func (handler *OperatorHandler) CreateOperator(ctx *gin.Context) {
	// Attempt to extract the body
	var newOperator model.Operator
	err := ctx.BindJSON(&newOperator)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Set the organization ID of the operator
	organization, _ := middleware.GetOrganizationClaim(ctx)
	newOperator.OrganizationId = organization.Id

	// Guard against non-admin users
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against non-unique operators in the organization
	if !handler.service.IsOperatorUnique(ctx, &newOperator) {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.OperatorNotUnique))
	}

	// Attempt to create the operator
	err = handler.service.Create(ctx.Request.Context(), &newOperator)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newOperator, "Successfully created operator.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OperatorHandler) GetOperators(ctx *gin.Context) {
	// Attempt to get the operators
	organization, _ := middleware.GetOrganizationClaim(ctx)
	operators, err := handler.service.FindByOrganizationId(ctx.Request.Context(), organization.Id)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Send the response
	result := utils.SuccessPayload(operators, "Successfully retrieved operators.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OperatorHandler) UpdateOperator(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedOperator model.Operator
	err := ctx.BindJSON(&updatedOperator)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin users
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to find the existing operator
	organization, _ := middleware.GetOrganizationClaim(ctx)
	operator, err := handler.service.FindById(ctx, updatedOperator.Id)
	if err != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.OperatorNotFound))
		return
	}

	// Guard against cross-tenant writing
	if organization.Id != operator.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against non-unique operator names
	updatedOperator.OrganizationId = operator.OrganizationId
	if !handler.service.IsOperatorUnique(ctx, &updatedOperator) {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.OperatorNotUnique))
		return
	}

	// Attempt to update the operator
	updatedOperator.OrganizationId = organization.Id
	err = handler.service.Update(ctx.Request.Context(), &updatedOperator)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated operator.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *OperatorHandler) DeleteOperator(ctx *gin.Context) {
	// Attempt to parse the query param
	organization, _ := middleware.GetOrganizationClaim(ctx)
	operatorIdToDelete, err := uuid.Parse(ctx.Param("operatorId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.OperatorNotFound))
		return
	}

	// Attempt to find the operator to delete
	operator, err := handler.service.FindById(ctx, operatorIdToDelete)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant deletions
	if organization.Id != operator.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the operator
	err = handler.service.Delete(ctx.Request.Context(), operator.Id)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted operator.")
	utils.Response(ctx, http.StatusOK, result)
}
