package controllers

import (
	"database-ms/databases"
	"encoding/json"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type UserFirestore struct {
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
	FireStore					*UserFirestore		`json:"firestore"`
}

func GetUsers(c *gin.Context) {
	organizationId := c.Param("organizationId")

	users := []User{}

	// Get all the user records
	iter := databases.Database.Client.
		Collection("organizations").
			Doc(organizationId).
				Collection("users").
					Documents(databases.Database.Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		
		user, authError := databases.Database.Auth.GetUser(databases.Database.Context, doc.Ref.ID)
		if authError != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		// Append to the list
		role := user.CustomClaims["role"].(string)
		users = append(users, User {
			UserId: &doc.Ref.ID,
			OrganizationId: &organizationId, 
			DisplayName: &user.DisplayName,
			Role: &role, 
			Email: &user.Email,
			Disabled: &user.Disabled,
			// Preferences not needed for this request
		})
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	organizationId := c.Param("organizationId")
	userId := c.Param("userId")

	// Get the data from authentication
	user, err := databases.Database.Auth.GetUser(databases.Database.Context, userId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get the firestore user data
	snapshot, err := databases.Database.Client. 
		Collection("organizations"). 
			Doc(organizationId).
				Collection("users").
					Doc(userId).
						Get(databases.Database.Context) 
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	firestore := UserFirestore{}
	data, parseError := json.Marshal(snapshot.Data())
	if parseError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	json.Unmarshal(data, &firestore)

	// Create the response
	role := user.CustomClaims["role"].(string)
	response := User {
		UserId: &userId,
		OrganizationId: &organizationId,
		DisplayName: &user.DisplayName,
		Role: &role,
		Email: &user.Email,
		Disabled: &user.Disabled,
		FireStore: &firestore,
	}

	c.JSON(http.StatusOK, response)
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
		c.Status(http.StatusInternalServerError)
	}

	// Create the user
	newUserParams := (&auth.UserToCreate{}).
		DisplayName(*newUser.DisplayName).
		Email(*newUser.Email).
		Password(*newUser.Password).
		Disabled(true)
	user, createError := databases.Database.Auth.CreateUser(databases.Database.Context, newUserParams)
	if createError != nil {
		// TODO: Send message indicating error, could be that email already taken, display name taken, etc
		c.Status(http.StatusInternalServerError)
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
		c.Status(http.StatusInternalServerError)
		return
	}

	// Create a user record for the organization
	_, fetchError := databases.Database.Client.
		Collection("organizations").
			Doc(*newUser.OrganizationID).
				Collection("users").
					Doc(user.UID).
						Set(databases.Database.Context, map[string]interface{} {
							"streamingSensors": []string{},
							// TODO: Probably many other settings in the future
						})
	if fetchError != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
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
		// TODO: Send error message if display name is taken
		c.Status(http.StatusOK)
		return
	}

	// Update the user's preferences
	_, fetchError := databases.Database.Client.
		Collection("organizations").
			Doc(*updatedUser.OrganizationId).
				Collection("users").
					Doc(*updatedUser.OrganizationId).
						Set(databases.Database.Context, *updatedUser.FireStore)
	if fetchError != nil {
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}