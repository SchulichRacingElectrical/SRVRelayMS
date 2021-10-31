package main

import (
	"database-ms/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Db config.FirebaseDB

func main() {
	Db = config.FirebaseDB{}
	err := Db.FirestoreDBInit()

	if err != nil {
		fmt.Println("Could not connect to db:\n", err.Error())
		return
	}

	router := gin.Default()
	databaseHandlers := router.Group("/database")
	{
		databaseHandlers.GET("/sensors", getSensors)
		databaseHandlers.GET("/sensors/:sid", getSensor)
		//PUT
		//DELETE
		//POST
	}

	router.Run(":8080")
}

// TODO: move these to follow industry project layout

func getSensors(c *gin.Context) {
	dsnap := Db.Client.Collection("sensors")

	// TODO: check for last update queryparam
	iter := dsnap.Documents(Db.Context)
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

func getSensor(c *gin.Context) {
	sid := c.Param("sid")

	dsnap, err := Db.Client.Collection("sensors").Doc(sid).Get(Db.Context)
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
