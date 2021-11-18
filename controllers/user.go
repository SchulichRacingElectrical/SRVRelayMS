package controllers

import (
	"database-ms/databases"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type User struct {

}

func GetUsers(c *gin.Context) {
	organization_id := c.Param("organization_id")

	type user_struct struct {
		user_id string
		display_name string
		role string
	}

	users := []user_struct{}

	// Get all the user records
	iter := databases.Database.Client.Collection("users").Where("organization_id", "==", organization_id).Documents(databases.Database.Context)
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
		
		user, authError := databases.Database.Auth.GetUser(databases.Database.Context, doc.Ref.ID)
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

func GetUser(c *gin.Context) {
	user_id := c.Param("user_id")

	snapshot, err := databases.Database.Client.Collection("users").Doc(user_id).Get(databases.Database.Context) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "User not found",
			"error": true,
		})
	}

	user := snapshot.Data()
	c.JSON(http.StatusOK, user)
}

func PostUser(c *gin.Context) {
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
	_, fetchError := databases.Database.Client.Collection("users").Doc(user.user_id).Set(databases.Database.Context, map[string]interface{}{
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

func PutUser(c *gin.Context) {
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
	_, fetchError := databases.Database.Client.Collection("users").Doc(user.user_id).Set(databases.Database.Context, c.Request.Body)
	if fetchError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message": "Could not update user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}