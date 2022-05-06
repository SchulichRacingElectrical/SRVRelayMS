package handlers

import (
	"database-ms/app/services"

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

}

func (handler *ChartPresetHandler) GetChartPresets(ctx *gin.Context) {

}

func (handler *ChartPresetHandler) UpdateChartPreset(ctx *gin.Context) {

}

func (handler *ChartPresetHandler) DeleteChartPreset(ctx *gin.Context) {
	
}