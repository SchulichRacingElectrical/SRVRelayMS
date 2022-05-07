package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/models"
	"database-ms/app/services"
	"database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RawDataPresetHandler struct {
	service				services.RawDataPresetServiceInterface
	thingService	services.ThingServiceInterface
}

func NewRawDataPresetAPI(service services.RawDataPresetServiceInterface, thingService services.ThingServiceInterface) *RawDataPresetHandler {
	return &RawDataPresetHandler{service: service, thingService: thingService}
}

func (handler *RawDataPresetHandler) CreateRawDataPreset(ctx *gin.Context) {
	var newRawDataPreset models.RawDataPreset
	ctx.BindJSON(&newRawDataPreset)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newRawDataPreset.ThingId.Hex())
	if err == nil {
		if handler.service.DoPresetSensorsExist(ctx, &newRawDataPreset) {
			if handler.service.IsRawDataPresetUnique(ctx, &newRawDataPreset) {
				if thing.OrganizationId == organization.ID {
					err := handler.service.Create(ctx.Request.Context(), &newRawDataPreset)
					if err == nil {
						result := utils.SuccessPayload(newRawDataPreset, "Successfully created Raw Data Preset.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
					}
				} else {
					utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
				}
			} else {
				utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.RawDataPresetNotUnique))
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))	
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
	}
}

func (handler *RawDataPresetHandler) GetRawDataPresets(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, ctx.Param("thingId"))
	if err == nil {
		if thing.OrganizationId == organization.ID {
			rawDataPresets, err := handler.service.FindByThingId(ctx.Request.Context(), ctx.Param("thingId"))
			if err == nil {
				result := utils.SuccessPayload(rawDataPresets, "Successfully retrieved Raw Data Presets.")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
	}
}

func (handler *RawDataPresetHandler) UpdateRawDataPreset(ctx *gin.Context) {
	var updatedRawDataPreset models.RawDataPreset
	ctx.BindJSON(&updatedRawDataPreset)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updatedRawDataPreset.ThingId.Hex())
	if err == nil {
		if handler.service.DoPresetSensorsExist(ctx, &updatedRawDataPreset) {
			if handler.service.IsRawDataPresetUnique(ctx, &updatedRawDataPreset) {
				if thing.OrganizationId == organization.ID {
					err := handler.service.Update(ctx.Request.Context(), &updatedRawDataPreset)
					if err == nil {
						result := utils.SuccessPayload(nil, "Successfully updated.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
					}
				} else {
					utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
				}
			} else {
				utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.RawDataPresetNotUnique))
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))	
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))	
	}
}

func (handler *RawDataPresetHandler) DeleteRawDataPreset(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	rawDataPreset, err := handler.service.FindById(ctx, ctx.Param("rpId"))
	if err == nil {
		thing, err := handler.thingService.FindById(ctx, rawDataPreset.ThingId.Hex())
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.service.Delete(ctx, ctx.Param("rpId"))
				if err == nil {
					result := utils.SuccessPayload(nil, "Successfully deleted.")
					utils.Response(ctx, http.StatusOK, result)
				} else {
					utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
				}
			} else {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))	
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.RawDataPresetNotFound))
	}
}

