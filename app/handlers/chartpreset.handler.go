package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/models"
	"database-ms/app/services"
	"database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChartPresetHandler struct {
	service 			services.ChartPresetServiceInterface
	thingService 	services.ThingServiceInterface
}

func NewChartPresetAPI(service services.ChartPresetServiceInterface, thingService services.ThingServiceInterface) *ChartPresetHandler {
	return &ChartPresetHandler{service: service, thingService: thingService}
}

func (handler *ChartPresetHandler) CreateChartPreset(ctx *gin.Context) {
	var newChartPreset models.ChartPreset
	ctx.BindJSON(&newChartPreset)
	if len(newChartPreset.Charts) == 0 {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotValid))	
		return
	}
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newChartPreset.ThingId.Hex())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			if handler.service.IsPresetValid(ctx, &newChartPreset) {
				if handler.service.IsPresetUnique(ctx, &newChartPreset) {
					err := handler.service.Create(ctx.Request.Context(), &newChartPreset)
					if err == nil {
						result := utils.SuccessPayload(newChartPreset, "Successfully created Chart Preset.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
					}
				} else {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotUnique))
				}
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotValid))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
	}
}

func (handler *ChartPresetHandler) GetChartPresets(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, ctx.Param("thingId"))
	if err == nil {
		if thing.OrganizationId == organization.ID {
			chartPresets, err := handler.service.FindByThingId(ctx.Request.Context(), ctx.Param("thingId"))
			if err == nil {
				result := utils.SuccessPayload(chartPresets, "Successfully retrieved Chart Presets.")
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

func (handler *ChartPresetHandler) UpdateChartPreset(ctx *gin.Context) {
	var updatedChartPreset models.ChartPreset
	ctx.BindJSON(&updatedChartPreset)
	if len(updatedChartPreset.Charts) == 0 {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotValid))	
		return
	}
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updatedChartPreset.ThingId.Hex())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			if handler.service.IsPresetValid(ctx, &updatedChartPreset) {
				if handler.service.IsPresetUnique(ctx, &updatedChartPreset) {
					err := handler.service.Update(ctx, &updatedChartPreset)
					if err == nil {
						result := utils.SuccessPayload(updatedChartPreset, "Successfully Updated.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
					}
				} else {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotUnique))
				}
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotValid))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))	
	}
}

func (handler *ChartPresetHandler) DeleteChartPreset(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	chartPreset, err := handler.service.FindById(ctx, ctx.Param("chartPresetId"))
	if err == nil {
		thing, err := handler.thingService.FindById(ctx, chartPreset.ThingId.Hex())
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.service.Delete(ctx, ctx.Param("chartPresetId"))
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
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ChartPresetNotFound))
	}
}