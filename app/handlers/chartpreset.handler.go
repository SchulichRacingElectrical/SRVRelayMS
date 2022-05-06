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
	ctx.BindJSON(newChartPreset)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newChartPreset.ThingId.Hex())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			if handler.service.AreChartsValid(ctx, &newChartPreset) {
				if handler.service.IsChartPresetUnique(ctx, &newChartPreset) {
					err := handler.service.Create(ctx.Request.Context(), &newChartPreset)
					if err == nil {
						result := utils.SuccessPayload(newChartPreset, "Successfully created Chart Preset.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						// Internal error
					}
				} else {
					// Bad request
				}
			} else {
				// Bad request
			}
		} else {
			// Bad Auth
		}
	} else {
		// Bad request
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
				// Internal error
			}
		} else {
			// No auth
		}
	} else {
		// Bad request
	}
}

func (handler *ChartPresetHandler) UpdateChartPreset(ctx *gin.Context) {
	var updatedChartPreset models.ChartPreset
	ctx.BindJSON(updatedChartPreset)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updatedChartPreset.ThingId.Hex())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			if handler.service.AreChartsValid(ctx, &updatedChartPreset) {
				if handler.service.IsChartPresetUnique(ctx, &updatedChartPreset) {
					err := handler.service.Update(ctx, &updatedChartPreset)
					if err == nil {
						result := utils.SuccessPayload(nil, "Successfully Updated.")
						utils.Response(ctx, http.StatusOK, result)
					} else {
						// Something went wrong
					}
				} else {
					// Bad request
				}
			} else {
				// Bad request
			}
		} else {
			// No auth
		}
	} else {
		// Bad request
	}
}

func (handler *ChartPresetHandler) DeleteChartPreset(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	chartPreset, err := handler.service.FindById(ctx, ctx.Param("cpId"))
	if err == nil {
		thing, err := handler.thingService.FindById(ctx, chartPreset.ThingId.Hex())
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.service.Delete(ctx, ctx.Param("cpId"))
				if err == nil {
					result := utils.SuccessPayload(nil, "Successfully deleted.")
					utils.Response(ctx, http.StatusOK, result)
				} else {
					// Internal error
				}
			} else {
				// Not auth
			}
		} else {
			// Bad request
		}
	} else {
		// Bad request
	}
}

// TODO: Indicate bad request rather than something else if a thing is found that does not belong to the claim