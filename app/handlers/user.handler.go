package handlers

import (
	"database-ms/app/models"
	userSrv "database-ms/app/services/user"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	user userSrv.UserServiceInterface
}

func NewUserAPI(userService userSrv.UserServiceInterface) *UserHandler {
	return &UserHandler{
		user: userService,
	}
}

func (handler *UserHandler) Create(c *gin.Context) {
	var newUser models.User
	c.BindJSON(&newUser)
	result := make(map[string]interface{})
	var status int

	// Check if this email already exists
	_, err := handler.user.FindByUserEmail(c.Request.Context(), newUser.Email)
	if err == nil {
		result = utils.NewHTTPError(utils.UserAlreadyExists)
		status = http.StatusConflict
		utils.Response(c, status, result)
		return
	}

	newUser.Password = hashPassword(newUser.Password)
	newUser.Roles = "Pending"
	res, err := handler.user.Create(c.Request.Context(), &newUser)
	if err == nil {
		result = utils.SuccessPayload(res, "Successfully created user")
		status = http.StatusOK
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	result := make(map[string]interface{})
	user, err := handler.user.FindByUserId(c.Request.Context(), c.Param("userId"))
	if err == nil {
		user.Password = ""
		result = utils.SuccessPayload(user, "Successfully retrieved user")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.UserNotFound)
		utils.Response(c, http.StatusBadRequest, result)
	}
}

func (handler *UserHandler) Login(c *gin.Context) {
	var loggingInUser models.User
	c.BindJSON(&loggingInUser)
	result := make(map[string]interface{})

	DBuser, err := handler.user.FindByUserEmail(c.Request.Context(), loggingInUser.Email)
	if err != nil {
		result = utils.NewHTTPError(utils.UserNotFound)
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	// Check password match
	if checkPasswordHash(loggingInUser.Password, DBuser.Password) {
		token, err := handler.user.CreateToken(c, DBuser)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.JSON(http.StatusOK, token)
	} else {
		result = utils.NewHTTPError(utils.WrongPassword)
		utils.Response(c, http.StatusUnauthorized, result)
	}
}

// Password hashing and verification functions

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic("Hashing password failed")
	}
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
