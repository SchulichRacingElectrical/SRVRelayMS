package controllers

import (
	"context"
	"database-ms/databases"
	"database-ms/models"
	"log"
	"net/http"
	"strconv"
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
	loc, err := time.LoadLocation("") // set timezome as UTC
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	sensor.LastUpdate = primitive.NewDateTimeFromTime(time.Now().In(loc))

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
