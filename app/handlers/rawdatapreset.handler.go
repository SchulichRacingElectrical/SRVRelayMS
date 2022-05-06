package handlers

import (
	"database-ms/app/services"

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

}

func (handler *RawDataPresetHandler) GetRawDataPresets(ctx *gin.Context) {

}

func (handler *RawDataPresetHandler) UpdateRawDataPreset(ctx *gin.Context) {

}

func (handler *RawDataPresetHandler) DeleteRawDataPreset(ctx *gin.Context) {
	
}

