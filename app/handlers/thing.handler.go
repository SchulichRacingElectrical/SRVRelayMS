package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/app/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ThingHandler struct {
	service  services.ThingServiceInterface
	filepath string
}

func NewThingAPI(thingService services.ThingServiceInterface, filepath string) *ThingHandler {
	return &ThingHandler{service: thingService, filepath: filepath}
}

func (handler *ThingHandler) CreateThing(ctx *gin.Context) {
	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to extract the body
	var newThing model.Thing
	err := ctx.BindJSON(&newThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to create the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	newThing.OrganizationId = organization.Id
	perr := handler.service.Create(ctx.Request.Context(), &newThing)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.ThingNotUnique))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(newThing, "Successfully created thing.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) GetThings(ctx *gin.Context) {
	// Attempt to fetch the things
	organization, _ := middleware.GetOrganizationClaim(ctx)
	things, perr := handler.service.FindByOrganizationId(ctx.Request.Context(), organization.Id)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(things, "Successfully retrieved things.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) UpdateThing(ctx *gin.Context) {
	// Guard against non-admin users
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to extract the body
	var updatedThing model.Thing
	err := ctx.BindJSON(&updatedThing)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the thing we are updated
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.service.FindById(ctx, updatedThing.Id)
	if perr != nil {
		utils.Response(ctx, http.StatusNotFound, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant writes
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the thing
	updatedThing.OrganizationId = organization.Id
	perr = handler.service.Update(ctx.Request.Context(), &updatedThing)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.ThingNotUnique))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Succesfully updated thing.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *ThingHandler) DeleteThing(ctx *gin.Context) {
	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Attempt to parse the query param
	thingIdToDelete, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.service.FindById(ctx, thingIdToDelete)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant deletion
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the thing
	perr = handler.service.Delete(ctx.Request.Context(), thingIdToDelete)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to delete session files related to the thing
	if err = os.RemoveAll(handler.filepath + thingIdToDelete.String()); err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.FailedToDeleteFiles))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted thing.")
	utils.Response(ctx, http.StatusOK, result)
}
