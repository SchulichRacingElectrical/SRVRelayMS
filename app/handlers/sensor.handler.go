package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SensorHandler struct {
	sensorService services.SensorServiceInterface
	thingService  services.ThingServiceInterface
}

func NewSensorAPI(sensorService services.SensorServiceInterface, thingService services.ThingServiceInterface) *SensorHandler {
	return &SensorHandler{sensorService: sensorService, thingService: thingService}
}

func (handler *SensorHandler) CreateSensor(ctx *gin.Context) {
	// Attempt to extract the body
	var newSensor model.Sensor
	err := ctx.BindJSON(&newSensor)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, newSensor.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant writing
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the sensor
	perr = handler.sensorService.Create(ctx.Request.Context(), &newSensor)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(newSensor, "Successfully created sensor")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SensorHandler) FindThingSensors(ctx *gin.Context) {
	// Attempt to read from the params
	thingId, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the sensors
	sensors, perr := handler.sensorService.FindByThingId(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(sensors, "Successfully retrieved sensors")
	utils.Response(ctx, http.StatusOK, result)
}

// TODO: Returns list of all sensors Ids
func (handler *SensorHandler) FindUpdatedSensors(ctx *gin.Context) {
	// Attempt to read from the params
	thingId, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to parse the last update
	lastUpdate, err := strconv.ParseInt(ctx.Param("lastUpdate"), 10, 64)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to fetch the updated sensors
	sensors, perr := handler.sensorService.FindUpdatedSensors(ctx.Request.Context(), thingId, lastUpdate)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(sensors, "Successfully retrieved sensors")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SensorHandler) UpdateSensor(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedSensor model.Sensor
	err := ctx.BindJSON(&updatedSensor)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to get the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, updatedSensor.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant updates
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the sensor
	perr = handler.sensorService.Update(ctx.Request.Context(), &updatedSensor)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.SensorNotUnique))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *SensorHandler) DeleteSensor(ctx *gin.Context) {
	// Attempt to read from the params
	sensorId, err := uuid.Parse(ctx.Param("sensorId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the sensor
	organization, _ := middleware.GetOrganizationClaim(ctx)
	sensor, perr := handler.sensorService.FindBySensorId(ctx, sensorId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, sensor.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant deletion
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the sensor
	perr = handler.sensorService.Delete(ctx.Request.Context(), sensorId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}
