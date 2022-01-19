package handlers

import (
	"database-ms/app/models"
	sensorSrv "database-ms/app/services/sensor"
	"database-ms/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	sensor sensorSrv.SensorServiceInterface
}

func NewSensorAPI(sensorService sensorSrv.SensorServiceInterface) *SensorHandler {
	return &SensorHandler{
		sensor: sensorService,
	}
}

func (handler *SensorHandler) Create(c *gin.Context) {
	var newSensor models.Sensor
	c.BindJSON(&newSensor)
	result := make(map[string]interface{})

	err := handler.sensor.Create(c.Request.Context(), &newSensor)
	var status int
	if err == nil {
		res := &createSensorRes{
			ID: newSensor.ID,
		}
		result = utils.SuccessPayload(res, "Succesfully created sensor")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *SensorHandler) FindThingSensors(c *gin.Context) {
	result := make(map[string]interface{})
	sensors, err := handler.sensor.FindByThingId(c.Request.Context(), c.Param("thingId"))
	if err == nil {
		result = utils.SuccessPayload(sensors, "Succesfully retrieved sensors")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SensorHandler) FindBySensorId(c *gin.Context) {
	result := make(map[string]interface{})
	sensor, err := handler.sensor.FindBySensorId(c.Request.Context(), c.Param("sensorId"))
	if err == nil {
		result = utils.SuccessPayload(sensor, "Succesfully retrieved sensor")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SensorHandler) FindUpdatedSensor(c *gin.Context) {
	result := make(map[string]interface{})

	lastUpdate, err := strconv.ParseInt(c.Param("lastUpdate"), 10, 64)
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	sensors, err := handler.sensor.FindUpdatedSensor(c.Request.Context(), c.Param("thingId"), lastUpdate)

	if err == nil {
		result = utils.SuccessPayload(sensors, "Succesfully retrieved sensors")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *SensorHandler) Update(c *gin.Context) {
	var updateSensor models.SensorUpdate
	c.BindJSON(&updateSensor)
	result := make(map[string]interface{})
	err := handler.sensor.Update(c.Request.Context(), c.Param("sensorId"), &updateSensor)
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(c, http.StatusOK, result)
}

func (handler *SensorHandler) Delete(c *gin.Context) {
	result := make(map[string]interface{})
	err := handler.sensor.Delete(c.Request.Context(), c.Param("sensorId"))
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(c, http.StatusOK, result)
}
