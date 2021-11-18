package controllers

import (
	"database-ms/databases"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sensor struct {
	Sid         *int    `json:"sid"`
	Type        *string `json:"type,omitempty"`
	LastUpdated *int    `json:"lastUpdated,omitempty"`
	Group       *string `json:"group,omitempty"`
	Category    *string `json:"category,omitempty"`
	Name        *string `json:"name,omitempty"`
	Frequency   *int    `json:"frequency,omitempty"`
	Unit        *string `json:"unit,omitempty"`
	CanId       *string `json:"canId,omitempty"`    //TODO: Comes in as hex but should be converted to longlong
	Disabled    *bool   `json:"disabled,omitempty"` //TODO: Should have default when empty
}

func PostSensor(c *gin.Context) {
	dsnap := databases.Database.Client.Collection("sensors")

	var newSensor Sensor
	if err := c.BindJSON(&newSensor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	//TODO: required fields for a sensor; should create custom Unmarshall that checks for required fields
	//TODO: verify how sid will be generated
	if newSensor.Sid == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "sid is required",
			"error":   true,
		})
		return
	}

	_, err := dsnap.Doc(strconv.Itoa(*newSensor.Sid)).Create(databases.Database.Context, newSensor)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Sensor with sid %d already exists", *newSensor.Sid),
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

	c.JSON(http.StatusOK, gin.H{
		"messaage": fmt.Sprintf("String with sid %d successfully created", *newSensor.Sid),
	})

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
		"message": fmt.Sprintf("String with sid %s successfully deleted", sid),
	})
}
