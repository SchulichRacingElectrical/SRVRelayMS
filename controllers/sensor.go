package controllers

import (
	"database-ms/databases"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddSensor
func AddSensor(c *gin.Context) {

}

func GetSensors(c *gin.Context) {

	dsnap := databases.Database.Client.Collection("sensors")

	// TODO: check for last update queryparam
	iter := dsnap.Documents(databases.Database.Context)
	sensors := make([]interface{}, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		sensors = append(sensors, doc.Data())
	}

	c.JSON(http.StatusOK, gin.H{
		"sensors": sensors,
	})

}

func GetSensor(c *gin.Context) {
	sid := c.Param("sid")

	dsnap, err := databases.Database.Client.Collection("sensors").Doc(sid).Get(databases.Database.Context)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sensor not found",
				"error":   true,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   true,
			})
		}
		return
	}

	sensor := dsnap.Data()
	c.JSON(http.StatusOK, sensor)
}

//UpdateSensor
func UpdateSensor(c *gin.Context) {

}

// DeleteSensor deletes a sensor document given the sid, if sensor does not exist then no error
func DeleteSensor(c *gin.Context) {
	sid := c.Param("sid")

	_, err := databases.Database.Client.Collection("sensors").Doc(sid).Delete(databases.Database.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messaage": "Sensor with sid = " + sid + " was successfully deleted",
	})
}
