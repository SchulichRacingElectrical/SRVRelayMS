package handlers

import (
	"database-ms/app/models"
	sensorSrv "database-ms/app/services/sensor"
	"database-ms/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type SensorHandler struct {
	snsr sensorSrv.SensorServiceInterface
}

func NewSensorAPI(sensorService sensorSrv.SensorServiceInterface) *SensorHandler {
	return &SensorHandler{
		snsr: sensorService,
	}
}

func (h *SensorHandler) Create(c *gin.Context) {
	var newSensor models.Sensor
	c.BindJSON(&newSensor)
	result := make(map[string]interface{})

	// Check if sensor already exists
	if h.snsr.IsSensorAlreadyExists(c.Request.Context(), newSensor.ThingID, newSensor.SID) {
		result = utils.NewHTTPError(utils.SensorAlreadyExists)
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	// TODO Check if Thing exists

	// Generate last updated
	newSensor.LastUpdate = utils.CurrentTimeInMilli()

	// Generate id
	newSensor.ID = bson.NewObjectId()

	err := h.snsr.Create(c.Request.Context(), &newSensor)
	var status int
	if err != nil {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	} else {
		res := &createSensorRes{
			ID: newSensor.ID,
		}
		result = utils.SuccessPayload(res, "Succesfully created sensor")
		status = http.StatusOK
	}
	utils.Response(c, status, result)
}

func (h *SensorHandler) FindByThingId(c *gin.Context) {

	thingId := c.Param("thingId")
	sensors, err := h.snsr.FindByThingId(c.Request.Context(), thingId)

	result := make(map[string]interface{})

	if err != nil {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	} else {
		result = utils.SuccessPayload(sensors, "Succesfully retrieved sensors")
		utils.Response(c, http.StatusOK, result)
	}

}

func (h *SensorHandler) FindById(c *gin.Context) {

	id := c.Param("id")
	sensor, err := h.snsr.FindById(c.Request.Context(), id)

	result := make(map[string]interface{})

	if err != nil {
		result = utils.NewHTTPError(utils.SensorNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	} else {
		result = utils.SuccessPayload(sensor, "Succesfully retrieved sensor")
		utils.Response(c, http.StatusOK, result)
	}

}

func (h *SensorHandler) FindByThingIdAndSid(c *gin.Context) {

	result := make(map[string]interface{})

	thingId := c.Param("thingId")
	sidStr := c.Param("sid")
	sid, err := strconv.Atoi(sidStr)
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	sensor, err := h.snsr.FindByThingIdAndSid(c.Request.Context(), thingId, sid)

	if err != nil {
		result = utils.NewHTTPError(utils.SensorNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	} else {
		result = utils.SuccessPayload(sensor, "Succesfully retrieved sensor")
		utils.Response(c, http.StatusOK, result)
	}

}

func (h *SensorHandler) FindByThingIdAndLastUpdate(c *gin.Context) {

	result := make(map[string]interface{})

	thingId := c.Param("thingId")
	lastUpdateStr := c.Param("lastUpdate")
	lastUpdate, err := strconv.ParseInt(lastUpdateStr, 10, 64)
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	sensors, err := h.snsr.FindByThingIdAndLastUpdate(c.Request.Context(), thingId, lastUpdate)

	if err != nil {
		result = utils.NewHTTPError(utils.SensorsNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	} else {
		result = utils.SuccessPayload(sensors, "Succesfully retrieved sensors")
		utils.Response(c, http.StatusOK, result)
	}

}

func (h *SensorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var updateSensor models.SensorUpdate
	c.BindJSON(&updateSensor)
	result := make(map[string]interface{})
	err := h.snsr.Update(c.Request.Context(), id, &updateSensor)
	if err != nil {
		result = utils.NewHTTCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(c, http.StatusOK, result)

}
