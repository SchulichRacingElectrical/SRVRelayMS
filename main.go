package main

import (
	"database-ms/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var firebase config.Firebase

func main() {
	// Configure firebase
	firebase = config.Firebase{}
	err := firebase.Init()
	if err != nil {
		fmt.Println("Could not connect to db:\n", err.Error())
		return
	}

	// Configure routes
	router := gin.Default()
	databaseHandlers := router.Group("/database")
	{
		// GET
		databaseHandlers.GET("/users/:organization_id")
		databaseHandlers.GET("/users/:organization_id/:user_id", getUser)
		databaseHandlers.GET("/sensors", getSensors)
		databaseHandlers.GET("/sensors/:sid", getSensor)
		// POST
		databaseHandlers.POST("/users", postUser)
		// PUT
		databaseHandlers.PUT("/users", putUser) // Should there be a param?
	}

	router.Run(":8080")
}

// TODO: move these to follow industry project layout

func getSensors(c *gin.Context) {
	dsnap := firebase.Client.Collection("sensors")

	// TODO: check for last update queryparam
	iter := dsnap.Documents(firebase.Context)
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

	dsnap, err := firebase.Client.Collection("sensors").Doc(sid).Get(firebase.Context)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sensor not found",
				"error": true,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error": true,
			})
		}
		return
	}

	sensor := dsnap.Data()
	c.JSON(http.StatusOK, sensor)
}

func getUsers(c *gin.Context) {
	organization_id := c.Param("organization_id")

	type user_struct struct {
		user_id string
		display_name string
		role string
	}

	users := []user_struct{}

	// Get all the user records
	iter := firebase.Client.Collection("users").Where("organization_id", "==", organization_id).Documents(firebase.Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"message": "Could not get users.",
			})
		}
		
		user, authError := firebase.Auth.GetUser(firebase.Context, doc.Ref.ID)
		if authError != nil {
			c.JSON(http.StatusInternalServerError, gin.H {
				"message": "Could not get users.",
			})
		}

		// Concatenate user and auth into list, then return
		users = append(users, user_struct{
			doc.Ref.ID,
			user.DisplayName,
			doc.Data()["role"].(string),
		})
	}

	data,parseError := json.Marshal(users)
	if parseError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "Parsing error.",
		})
	}
	c.JSON(http.StatusOK, data)
}

func getUser(c *gin.Context) {
	user_id := c.Param("user_id")

	snapshot, err := firebase.Client.Collection("users").Doc(user_id).Get(firebase.Context) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "User not found",
			"error": true,
		})
	}

	user := snapshot.Data()
	c.JSON(http.StatusOK, user)
}

func postUser(c *gin.Context) {
	// Get body
	type user_struct struct {
		organization_id string
		user_id string
	}
	decoder := json.NewDecoder(c.Request.Body)	
	var user user_struct
	decodeError := decoder.Decode(&user)
	if decodeError != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message": "Body is incorrectly formatted.",
		})
		return
	}

	// Create the user entry
	_, fetchError := firebase.Client.Collection("users").Doc(user.user_id).Set(firebase.Context, map[string]interface{}{
		"organization_id": user.organization_id,
		"role": "member",
	})
	if fetchError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "Could not create entry",
			"error": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{}) // Don't want to return anything...
}

func putUser(c *gin.Context) {
	// Get the body
	type user_struct struct {
		organization_id string
		user_id string
		role string
	}
	decoder := json.NewDecoder(c.Request.Body)
	var user user_struct
	decodeError := decoder.Decode(&user)
	if decodeError != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message": "Body is incorrectly formatted.",
		})
	}

	// Update the user
	_, fetchError := firebase.Client.Collection("users").Doc(user.user_id).Set(firebase.Context, c.Request.Body)
	if fetchError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "Could not update user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
