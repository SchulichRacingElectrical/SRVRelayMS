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
	service 			services.SensorServiceInterface
	thingService 	services.ThingServiceInterface
}

func NewSensorAPI(service services.SensorServiceInterface, thingService services.ThingServiceInterface) *SensorHandler {
	return &SensorHandler{service: service, thingService: thingService}
}

func (handler *SensorHandler) CreateSensor(ctx *gin.Context) {
	var newSensor models.Sensor
	err := ctx.BindJSON(&newSensor)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))	
		return	
	}
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, newSensor.ThingID.Hex())
	if handler.service.IsSensorUnique(ctx, &newSensor) {
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.service.Create(ctx.Request.Context(), &newSensor)
				if err == nil {
					result := utils.SuccessPayload(newSensor, "Successfully created Sensor.")
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
			sensors, err := handler.service.FindByThingId(ctx.Request.Context(), ctx.Param("thingId"))
			if err == nil {
				result := utils.SuccessPayload(sensors, "Successfully retrieved Sensors.")
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
				sensors, err := handler.service.FindUpdatedSensors(ctx.Request.Context(), ctx.Param("thingId"), lastUpdate)
				if err == nil {
					result := utils.SuccessPayload(sensors, "Successfully retrieved Sensors.")
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
	var updatedSensor models.Sensor
	ctx.BindJSON(&updatedSensor)
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, err := handler.thingService.FindById(ctx, updatedSensor.ThingID.Hex())
	if err == nil {
		if thing.OrganizationId == organization.ID {
			if handler.service.IsSensorUnique(ctx, &updatedSensor) {
				err := handler.service.Update(ctx.Request.Context(), &updatedSensor)
				if err == nil {
					result := utils.SuccessPayload(nil, "Successfully updated.")
					utils.Response(ctx, http.StatusOK, result)
				} else {
					utils.Response(ctx, http.StatusInternalServerError, utils.NewHTTPCustomError(utils.InternalError, err.Error()))
				}
			} else {
				utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(utils.SensorNotUnique))
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
	sensor, err := handler.service.FindBySensorId(ctx, ctx.Param("sensorId"))
	if err == nil {
		thing, err := handler.thingService.FindById(ctx, sensor.ThingID.Hex())	
		if err == nil {
			if thing.OrganizationId == organization.ID {
				err := handler.service.Delete(ctx.Request.Context(), ctx.Param("sensorId"))
				if err != nil {
					utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
				} else {
					result := utils.SuccessPayload(nil, "Successfully deleted.")
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
