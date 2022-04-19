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
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newSensor.ThingID.String())
	if handler.sensorService.IsSensorUnique(ctx, &newSensor) {
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
	} else {
		utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.SensorNotUnique))
	}
}

func (handler *SensorHandler) FindThingSensors(ctx *gin.Context) {
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
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, ctx.Param("thingId"))
	if err == nil {
		if thing.OrganizationId == organization.ID {
			lastUpdate, err := strconv.ParseInt(ctx.Param("lastUpdate"), 10, 64)
			if err != nil {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			} else {
				sensors, err := handler.sensorService.FindUpdatedSensors(ctx.Request.Context(), ctx.Param("thingId"), lastUpdate)
				if err == nil {
					result := utils.SuccessPayload(sensors, "Successfully retrieved sensors")
					utils.Response(ctx, http.StatusOK, result)
				} else {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
				}
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))		
	}
}

func (handler *SensorHandler) UpdateSensor(ctx *gin.Context) {
	var updateSensor models.Sensor
	ctx.BindJSON(&updateSensor)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updateSensor.ThingID.String())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			err := handler.sensorService.Update(ctx.Request.Context(), &updateSensor)
			if err != nil {
				utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
			} else {
				result := utils.SuccessPayload(nil, "Successfully updated")
				utils.Response(ctx, http.StatusOK, result)
			}
		} else {
			utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))		
	}
}

func (handler *SensorHandler) DeleteSensor(ctx *gin.Context) {
	organization, _ := middleware.GetOrganizationClaim(ctx)
	sensor, err := handler.sensorService.FindBySensorId(ctx, ctx.Param("sensorId"))
	if err == nil {
		thing, err := handler.thingService.FindById(ctx, sensor.ThingID.String())	
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.sensorService.Delete(ctx.Request.Context(), ctx.Param("sensorId"))
				if err != nil {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
				} else {
					result := utils.SuccessPayload(nil, "Successfully deleted")
					utils.Response(ctx, http.StatusOK, result)
				}
			} else {
				utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))	
			}
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))	
		}
	} else {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))		
	}
}

// Create tenant security function? - Might make more sense to put this in the service.