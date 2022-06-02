package handlers

import (
	"database-ms/app/services"

	"github.com/gin-gonic/gin"
)

type DatumHandler struct {
	datumService   services.DatumServiceInterface
	thingService   services.ThingServiceInterface
	sessionService services.SessionServiceInterface
}

func NewDatumAPI(
	datumService services.DatumServiceInterface,
	thingService services.ThingServiceInterface,
	sessionService services.SessionServiceInterface,
) *DatumHandler {
	return &DatumHandler{
		datumService:   datumService,
		thingService:   thingService,
		sessionService: sessionService,
	}
}

func (handler *DatumHandler) GetSensorData(ctx *gin.Context) {
	// TODO
}
