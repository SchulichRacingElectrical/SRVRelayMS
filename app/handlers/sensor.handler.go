package handlers

import (
	"database-ms/app/middleware"
	models "database-ms/app/models"
	services "database-ms/app/services"
	utils "database-ms/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	sensorService services.SensorServiceInterface
	thingService services.ThingServiceInterface
}

func NewSensorAPI(sensorService services.SensorServiceInterface, thingService services.ThingServiceInterface) *SensorHandler {
	return &SensorHandler{sensorService: sensorService, thingService: thingService}
}

func (handler *SensorHandler) CreateSensor(ctx *gin.Context) {
	var newSensor models.Sensor
	ctx.BindJSON(&newSensor)

	// Prevent cross-tenant creation of sensors
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newSensor.ThingID.String())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			err := handler.sensorService.Create(ctx.Request.Context(), &newSensor)
			if err == nil {
				result := utils.SuccessPayload(newSensor, "Successfully created sensor")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.EntityCreationError))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))	
	}
}

func (handler *SensorHandler) FindThingSensors(ctx *gin.Context) {
	// Prevent cross-tenant reading of sensors
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, ctx.Param("thingId"))
	if err == nil {
		if thing.OrganizationId == organization.ID {
			sensors, err := handler.sensorService.FindByThingId(ctx.Request.Context(), ctx.Param("thingId"))
			if err == nil {
				result := utils.SuccessPayload(sensors, "Successfully retrieved sensors")
				utils.Response(ctx, http.StatusOK, result)
			} else {
				utils.Response(ctx, http.StatusBadRequest,  utils.NewHTTPError(utils.SensorsNotFound))
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))		
	}
}

func (handler *SensorHandler) FindUpdatedSensors(ctx *gin.Context) {
	result := make(map[string]interface{})
	lastUpdate, err := strconv.ParseInt(ctx.Param("lastUpdate"), 10, 64)
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}
	sensors, err := handler.sensorService.FindUpdatedSensor(ctx.Request.Context(), ctx.Param("thingId"), lastUpdate)
	if err == nil {
		result = utils.SuccessPayload(sensors, "Successfully retrieved sensors")
		utils.Response(ctx, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(ctx, http.StatusBadRequest, result)
	}
}

// TODO: Ensure the thing for the updated sensor belongs to the correct organization
func (handler *SensorHandler) UpdateSensor(ctx *gin.Context) {
	var updateSensor models.SensorUpdate
	ctx.BindJSON(&updateSensor)
	result := make(map[string]interface{})
	err := handler.sensorService.Update(ctx.Request.Context(), ctx.Param("sensorId"), &updateSensor)
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}
	result = utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

// TODO: Ensure the thing for the deleted sensor belongs to the correct organization
func (handler *SensorHandler) DeleteSensor(ctx *gin.Context) {
	result := make(map[string]interface{})
	err := handler.sensorService.Delete(ctx.Request.Context(), ctx.Param("sensorId"))
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(ctx, http.StatusBadRequest, result)
		return
	}
	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}
