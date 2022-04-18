package handlers

import (
	"database-ms/app/models"
	thingSrv "database-ms/app/services/thing"
	"database-ms/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThingHandler struct {
	thing thingSrv.ThingServiceInterface
}

func NewThingAPI(thingService thingSrv.ThingServiceInterface) *ThingHandler {
	return &ThingHandler{
		thing: thingService,
	}
}

func (handler *ThingHandler) Create(c *gin.Context) {
	var newThing models.Thing
	c.BindJSON(&newThing)
	result := make(map[string]interface{})

	err := handler.thing.Create(c.Request.Context(), &newThing)
	var status int
	if err == nil {
		res := &createEntityRes{
			ID: newThing.ID,
		}
		result = utils.SuccessPayload(res, "Succesfully created thing")
	} else {
		fmt.Println(err)
		result = utils.NewHTTPError(utils.EntityCreationError)
		status = http.StatusBadRequest
	}
	utils.Response(c, status, result)

}

// Need to get ALL the things, not just by the ID
func (handler *ThingHandler) GetThings(c *gin.Context) {
	result := make(map[string]interface{})
	things, err := handler.thing.FindByOrganizationId(c.Request.Context(), c.Param("organizationId"))
	if err == nil {
		result = utils.SuccessPayload(things, "Successfully retrieved things.")
		utils.Response(c, http.StatusOK, result)
	} else {
		result = utils.NewHTTPError(utils.ThingNotFound)
		utils.Response(c, http.StatusBadRequest, result)	
	}
}

func (handler *ThingHandler) UpdateThing(c *gin.Context) {
	var thingUpdates models.ThingUpdate
	c.BindJSON(&thingUpdates)
	result := make(map[string]interface{})
	err := handler.thing.Update(c.Request.Context(), c.Param("thingId"), &thingUpdates)
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Succesfully updated")
	utils.Response(c, http.StatusOK, result)
}

func (handler *ThingHandler) Delete(c *gin.Context) {
	result := make(map[string]interface{})
	err := handler.thing.Delete(c.Request.Context(), c.Param("thingId"))
	if err != nil {
		result = utils.NewHTTPCustomError(utils.BadRequest, err.Error())
		utils.Response(c, http.StatusBadRequest, result)
		return
	}

	result = utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(c, http.StatusOK, result)
}
