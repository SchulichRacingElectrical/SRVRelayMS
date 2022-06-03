package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// Attempt to read from the params
	sessionId, err := uuid.Parse(ctx.Param("sessionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}
	sensorId, err := uuid.Parse(ctx.Param("sensorId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the session
	session, perr := handler.sessionService.FindById(ctx, sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.SessionNotFound))
		return
	}

	// If the session file is uploaded, its data is not available
	if *session.Generated != true {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the sensor
	sensor, perr := handler.sessionService.FindById(ctx, sensorId)
	if perr != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.SensorNotFound))
		return
	}

	// Attempt to find the session and sensor thing
	sessionThing, perr := handler.thingService.FindById(ctx, session.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.ThingNotFound))
		return
	}
	sensorThing, perr := handler.thingService.FindById(ctx, sensor.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against non-matching things between the session and sensor
	if sessionThing.Id != sensorThing.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against cross-tenant reads
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if sessionThing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to get the data
	data, perr := handler.datumService.FindBySessionIdAndSensorId(ctx, sessionId, sessionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.InternalError))
		return
	}

	// Send the response
	result := utils.SuccessPayload(data, "Successfully updated.")
	utils.Response(ctx, http.StatusOK, result)
}
