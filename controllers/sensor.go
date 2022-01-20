package controllers

import (
	"context"
	"database-ms/app/models"
	"database-ms/databases"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSensors(c *gin.Context) {
	// last updated query param should be handled
	filter := bson.M{}

	// handle _id query param
	if id := c.Query("id"); id != "" {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		filter["_id"] = objId
	}

	// handle thingId param
	if thingId := c.Query("thingId"); thingId != "" {
		objId, err := primitive.ObjectIDFromHex(thingId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"error":   true,
			})
			return
		}
		filter["thingId"] = objId

		// handle sid param
		// sid is only useful with thingid
		if sid := c.Query("sid"); sid != "" {
			i, err := strconv.Atoi(sid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "non integer sid",
					"error":   true,
				})
				return
			}
			filter["sid"] = i
		}
	}

	// TODO handle lastUpdate param here

	// query database
	collection := databases.Mongo.Db.Collection("Sensor")
	var sensors []models.Sensor
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "non integer sid",
			"error":   true,
		})
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var sensor models.Sensor
		err := cur.Decode(&sensor)
		if err != nil {
			log.Fatal(err)
		}

		sensors = append(sensors, sensor)
	}

	if err := cur.Err(); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  sensors,
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
