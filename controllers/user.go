package controllers

import (
	"database-ms/databases"
	utils "database-ms/utils"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type UserFirestore struct {
	StreamingSensors *[]string `json:"streamingSensors"`
}

type User struct {
	UserId         *string        `json:"userId"`
	OrganizationId *string        `json:"organizationId"`
	DisplayName    *string        `json:"displayName"`
	Role           *string        `json:"role"` // guest | admin | lead | member
	Email          *string        `json:"email"`
	Disabled       *bool          `json:"disabled"`
	FireStore      *UserFirestore `json:"firestore,omitempty"`
}

func GetUsers(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)

	users := []User{}
	iter := databases.Firebase.Client.
		Collection("organizations").
		Doc(organizationId).
		Collection("users").
		Documents(databases.Firebase.Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		user, authError := databases.Firebase.Auth.GetUser(databases.Firebase.Context, doc.Ref.ID)
		if authError != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		role := user.CustomClaims["role"].(string)
		users = append(users, User{
			UserId:         &doc.Ref.ID,
			OrganizationId: &organizationId,
			DisplayName:    &user.DisplayName,
			Role:           &role,
			Email:          &user.Email,
			Disabled:       &user.Disabled,
			// Preferences not needed for this request
		})
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	organization := c.GetStringMap("organization")
	organizationId := organization["organizationId"].(string)
	userId := c.Param("userId")

	// Get the data from authentication
	user, err := databases.Firebase.Auth.GetUser(databases.Firebase.Context, userId)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get the firestore user data
	snapshot, err := databases.Firebase.Client.
		Collection("organizations").
		Doc(organizationId).
		Collection("users").
		Doc(userId).
		Get(databases.Firebase.Context)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	firestore := UserFirestore{}
	utils.JsonToStruct(snapshot.Data(), &firestore)

	// Create the response
	role := user.CustomClaims["role"].(string)
	response := User{
		UserId:         &userId,
		OrganizationId: &organizationId,
		DisplayName:    &user.DisplayName,
		Role:           &role,
		Email:          &user.Email,
		Disabled:       &user.Disabled,
		FireStore:      &firestore,
	}
	c.JSON(http.StatusOK, response)
}

type PostUserBody struct {
	OrganizationID *string `json:"organizationId"`
	Email          *string `json:"email"`
	Password       *string `json:"password"`
	DisplayName    *string `json:"displayName"`
}

func PostUser(c *gin.Context) {
	var newUser PostUserBody
	if err := c.BindJSON(&newUser); err != nil {
		c.Status(http.StatusInternalServerError)
	}

	// Check if the organization exists
	snapshot := databases.Firebase.Client.Collection("organizations").Doc(*newUser.OrganizationID)
	organization, err := snapshot.Get(databases.Firebase.Context)
	if !organization.Exists() || err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	// Create the user
	newUserParams := (&auth.UserToCreate{}).
		DisplayName(*newUser.DisplayName).
		Email(*newUser.Email).
		Password(*newUser.Password).
		Disabled(true)
	user, createError := databases.Firebase.Auth.CreateUser(databases.Firebase.Context, newUserParams)
	if createError != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Create the custom claims for the user
	role := "member"
	users, err := snapshot.Collection("users").Documents(databases.Firebase.Context).GetAll()
	if len(users) == 0 || err != nil {
		role = "admin"
	}
	updateParams := (&auth.UserToUpdate{}).
		CustomClaims(map[string]interface{}{
			"role":           role,
			"organziationId": *newUser.OrganizationID,
		})
	_, authUpdateError := databases.Firebase.Auth.
		UpdateUser(databases.Firebase.Context, user.UID, updateParams)
	if authUpdateError != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Create a user record in the organization
	_, fetchError := databases.Firebase.Client.
		Collection("organizations").
		Doc(*newUser.OrganizationID).
		Collection("users").
		Doc(user.UID).
		Set(databases.Firebase.Context, map[string]interface{}{
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
		c.Status(http.StatusBadRequest)
	}

	// Update the auth data
	params := (&auth.UserToUpdate{}).
		DisplayName(*updatedUser.DisplayName).
		Email(*updatedUser.Email).
		Disabled(*updatedUser.Disabled).
		CustomClaims(map[string]interface{}{
			"role":           *updatedUser.Role,
			"organizationId": *updatedUser.OrganizationId,
		})
	_, authUpdateError := databases.Firebase.Auth.
		UpdateUser(databases.Firebase.Context, *updatedUser.UserId, params)
	if authUpdateError != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": authUpdateError.Error(),
		})
		return
	}

	// Update the user's preferences
	_, fetchError := databases.Firebase.Client.
		Collection("organizations").
		Doc(*updatedUser.OrganizationId).
		Collection("users").
		Doc(*updatedUser.OrganizationId).
		Set(databases.Firebase.Context, *updatedUser.FireStore)
	if fetchError != nil {
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}
