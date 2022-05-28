package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	"database-ms/app/services"
	"database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RawDataPresetHandler struct {
	service      services.RawDataPresetServiceInterface
	thingService services.ThingServiceInterface
}

func NewRawDataPresetAPI(service services.RawDataPresetServiceInterface, thingService services.ThingServiceInterface) *RawDataPresetHandler {
	return &RawDataPresetHandler{service: service, thingService: thingService}
}

func (handler *RawDataPresetHandler) CreateRawDataPreset(ctx *gin.Context) {
	// Attempt to parse the body
	var newRawDataPreset model.RawDataPreset
	err := ctx.BindJSON(&newRawDataPreset)
	if len(newRawDataPreset.SensorIds) == 0 || err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.RawDataPresetNotValid))
		return
	}

	// Attempt to find the associated thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newRawDataPreset.ThingId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant creation
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the preset
	err = handler.service.Create(ctx.Request.Context(), &newRawDataPreset)
	if err != nil {
		// handle error codes
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newRawDataPreset, "Successfully created Raw Data Preset.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *RawDataPresetHandler) GetRawDataPresets(ctx *gin.Context) {
	// Attempt to parse the params
	thingId, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the associated thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, thingId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the presets
	rawDataPresets, err := handler.service.FindByThingId(ctx.Request.Context(), ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(rawDataPresets, "Successfully retrieved Raw Data Presets.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *RawDataPresetHandler) UpdateRawDataPreset(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedRawDataPreset model.RawDataPreset
	err := ctx.BindJSON(&updatedRawDataPreset)
	if len(updatedRawDataPreset.SensorIds) == 0 || err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.RawDataPresetNotValid))
		return
	}

	// Attempt to read the associated thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updatedRawDataPreset.ThingId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant updates
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the preset
	err = handler.service.Update(ctx.Request.Context(), &updatedRawDataPreset)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *RawDataPresetHandler) DeleteRawDataPreset(ctx *gin.Context) {
	// Attempt to parse the params
	rawDataPresetId, err := uuid.Parse(ctx.Param("rawDataPresetId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the preset
	organization, _ := middleware.GetOrganizationClaim(ctx)
	rawDataPreset, err := handler.service.FindById(ctx, rawDataPresetId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.RawDataPresetNotFound))
		return
	}

	// Attempt to find the associated thing
	thing, err := handler.thingService.FindById(ctx, rawDataPreset.ThingId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant deletions
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt deletion
	err = handler.service.Delete(ctx, rawDataPresetId)
	if err != nil {
		utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted.")
	utils.Response(ctx, http.StatusOK, result)
}
