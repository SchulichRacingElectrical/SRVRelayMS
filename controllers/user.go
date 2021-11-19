package controllers

import (
	"database-ms/databases"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type UserPreferences struct {
	StreamingSensors	*[]string					`json:"streamingSensors"`
}

// TODO: Store everything in firebase's auth. Including custom claims
type User struct {
	UserId 						*string  					`json:"userId"`
	OrganizationId 		*string 					`json:"organizationId"`
	DisplayName 			*string						`json:"displayName"`
	Role 							*string						`json:"role"` // guest | admin | lead | member
	Email           	*string   				`json:"email"`
	Disabled					*bool 						`json:"disabled"`
	Preferences				*UserPreferences	`json:"preferences,omitempty"`
}

func GetUsers(c *gin.Context) {
	// organizationId := c.Param("organizationId")

	// users := []User{}

	// // Get all the user records
	// iter := databases.Database.Client.
	// 	Collection("organization").
	// 		Doc(organizationId).
	// 			Collection("users").
	// 				Documents(databases.Database.Context)
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H {
	// 			"message": "Could not get users.",
	// 		})
	// 	}
		
	// 	user, authError := databases.Database.Auth.GetUser(databases.Database.Context, doc.Ref.ID)
	// 	if authError != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H {
	// 			"message": "Could not get users.",
	// 		})
	// 	}

	// 	// Read custom claims
	// 	role := user.CustomClaims["role"].(string)

	// 	// Append to the list
	// 	users = append(users, User {
	// 		UserId: &doc.Ref.ID,
	// 		OrganizationId: &organizationId, 
	// 		DisplayName: &user.DisplayName,
	// 		Role: &role, 
	// 		Email: &user.Email,
	// 		Disabled: &user.Disabled,
	// 	})
	// }

	// c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	// organizationId := c.Param("organizationId")
	// userId := c.Param("userId")

	// // Get the data from firebase
	// _, err := databases.Database.Client. // TODO: change _ to snapshot
	// 	Collection("organization"). // TODO: Change firebase to have "organizations"
	// 		Doc(organizationId).
	// 			Collection("users").
	// 				Doc(userId).
	// 					Get(databases.Database.Context) 
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H {
	// 		"message": "User not found",
	// 		"error": true,
	// 	})
	// 	return
	// }
	// //userData = snapshot.Data()

	// // Get the data from authentication
	// auth, err := databases.Database.Auth.GetUser(databases.Database.Context, userId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H {
	// 		"message": "User not found",
	// 		"error": true,
	// 	})
	// 	return
	// }

	// role := auth.CustomClaims["role"].(string)
	// user := User {
	// 	UserId: &userId,
	// 	OrganizationId: &organizationId,
	// 	DisplayName: &auth.DisplayName,
	// 	Role: &role,
	// 	Email: &auth.Email,
	// 	Disabled: &auth.Disabled,
	// }

	// c.JSON(http.StatusOK, user)
}

func PostUser(c *gin.Context) {
	type PostUserBody struct {
		OrganizationID 	*string	`json:"organizationId"`
		Email						*string `json:"email"`
		Password				*string `json:"password"`
		DisplayName			*string `json:"displayName"`
	}

	// Parse the incoming body
	var newUser PostUserBody
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {})
	}

	// Create the user
	newUserParams := (&auth.UserToCreate{}).
		DisplayName(*newUser.DisplayName).
		Email(*newUser.Email).
		Password(*newUser.Password).
		Disabled(true)
	user, createError := databases.Database.Auth.CreateUser(databases.Database.Context, newUserParams)
	if createError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	// Create the custom claims for the user
	updateParams := (&auth.UserToUpdate{}).
		CustomClaims(map[string]interface{} {
			"role": "member", // New users will be members by default, TODO: Create an admin for the first user
		})
	_, authUpdateError := databases.Database.Auth.
		UpdateUser(databases.Database.Context, user.UID, updateParams)
	if authUpdateError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {})
		return
	}

	// Create a user record for the organization

	// Create the user entry
	// _, fetchError := databases.Database.Client.
	// 	Collection("organization").
	// 		Doc(*newUser.OrganizationID).
	// 			Collection("users").
	// 				Doc(*newUser.UserID).
	// 					Set(databases.Database.Context, map[string]interface{}{
	// 						"data": "x",
	// 					})
	// if fetchError != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H {
	// 		"message": "Could not set firebase document.",
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
	})
}

func PutUser(c *gin.Context) {
	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"message": "Body is incorrectly formatted.",
		})
	}

	// Update the auth data
	params := (&auth.UserToUpdate{}).
		DisplayName(*updatedUser.DisplayName).
		Email(*updatedUser.Email).
		Disabled(*updatedUser.Disabled).
		CustomClaims(map[string]interface{} {
			"role": *updatedUser.Role,
		})
	_, authUpdateError := databases.Database.Auth.
		UpdateUser(databases.Database.Context, *updatedUser.UserId, params)
	if authUpdateError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {})
		return
	}

	// Update the user's preferences
	_, fetchError := databases.Database.Client.
		Collection("organizations").
			Doc(*updatedUser.OrganizationId).
				Collection("users").
					Doc(*updatedUser.OrganizationId).
						Set(databases.Database.Context, *updatedUser.Preferences)
	if fetchError != nil {
		c.JSON(http.StatusInternalServerError, gin.H {})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}