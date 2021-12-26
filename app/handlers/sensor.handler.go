package handlers

import (
	"database-ms/app/models"
	sensorSrv "database-ms/app/services/sensor"
	"database-ms/utils"
	"fmt"
	"net/http"

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
