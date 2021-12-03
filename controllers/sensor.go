package controllers

import (
	"context"
	"database-ms/databases"
	"database-ms/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSensor(c *gin.Context) {
	// decode body request param
	var sensor models.Sensor
	if err := c.BindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	collection := databases.Mongo.Db.Collection("Sensor")

	// TODO autogenerate sid
	// validate sid is unique for all sensors in the thing
	filter := bson.M{"thingId": sensor.ThingID, "sid": sensor.SID}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "sid already exist",
			"error":   true,
		})
		return
	}

	// generate create time
	sensor.LastUpdate = primitive.NewDateTimeFromTime(time.Now())

	// write to database
	result, err := collection.InsertOne(context.TODO(), sensor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "created new sensor",
		"data":    result,
	})

}

func GetSensors(c *gin.Context) {
	// last updated query param should be handled
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "endpoint not available",
	})
}

func GetSensor(c *gin.Context) {
	// sid := c.Param("sid")
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "endpoint not available",
	})

}

func UpdateSensor(c *gin.Context) {
	// sid := c.Param("sid")
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "endpoint not available",
	})
}

func DeleteSensor(c *gin.Context) {
	// sid := c.Param("sid")
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "endpoint not available",
	})
}
