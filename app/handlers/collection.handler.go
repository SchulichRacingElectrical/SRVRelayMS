package handlers

import (
	"database-ms/app/middleware"
	"database-ms/app/model"
	"database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	collectionService services.CollectionServiceInterface
	thingService      services.ThingServiceInterface
}

func NewCollectionAPI(collectionService services.CollectionServiceInterface,
	thingService services.ThingServiceInterface) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
		thingService:      thingService,
	}
}

func (handler *CollectionHandler) CreateCollection(ctx *gin.Context) {
	var newCollection model.Collection
	err := ctx.BindJSON(&newCollection)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Admin") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Guard against cross-tenant writes
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, newCollection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingsNotFound))
		return
	}
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the collection
	perr = handler.collectionService.CreateCollection(ctx.Request.Context(), &newCollection)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(newCollection, "Successfully created collection.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) GetCollections(ctx *gin.Context) {
	// Attempt to read from the params
	thingId, err := uuid.Parse(ctx.Param("thingId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the collections
	collections, perr := handler.collectionService.GetCollectionsByThingId(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(collections, "Successfully retrieved collections")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) UpdateCollections(ctx *gin.Context) {
	// Attempt to extract the body
	var updatedCollection model.Collection
	err := ctx.BindJSON(&updatedCollection)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-admin lead or admin
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to get the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, updatedCollection.ThingId) // Note: FindById does not work when Thing DNE
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant updates
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to update the collection
	perr = handler.collectionService.UpdateCollection(ctx.Request.Context(), &updatedCollection)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError(perr.Error()))
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		}
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) DeleteCollection(ctx *gin.Context) {
	// Attempt to read from the params
	collectionId, err := uuid.Parse(ctx.Param("collectionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the collection
	organization, _ := middleware.GetOrganizationClaim(ctx)
	collection, perr := handler.collectionService.FindById(ctx.Request.Context(), collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SensorsNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant deletion
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the collection
	perr = handler.collectionService.DeleteCollection(ctx.Request.Context(), collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) AddComment(c *gin.Context) {
	// var comment model.Comment
	// c.BindJSON(&comment)

	// _, err := handler.session.FindById(c.Request.Context(), c.Param("sessionId"))
	// if err == nil {
	// 	err := handler.comment.AddComment(c.Request.Context(), utils.Session, c.Param("sessionId"), &comment)
	// 	if err == nil {
	// 		result := utils.SuccessPayload(nil, "Successfully added comment.")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		utils.Response(c, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
	// 	}
	// } else {
	// 	utils.Response(c, http.StatusNotFound, utils.NewHTTPError(utils.SessionNotFound))
	// }
}

func (handler *CollectionHandler) GetComments(c *gin.Context) {
	// comments, err := handler.comment.GetComments(c.Request.Context(), utils.Session, c.Param("sessionId"))
	// if err == nil {
	// 	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	utils.Response(c, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
	// }
}

func (handler *CollectionHandler) UpdateCommentContent(c *gin.Context) {
	// var updatedComment model.Comment
	// c.BindJSON(&updatedComment)

	// err := handler.comment.UpdateCommentContent(c.Request.Context(), c.Param("commentId"), &updatedComment)
	// if err == nil {
	// 	result := utils.SuccessPayload(nil, "Successfully updated comment")
	// 	utils.Response(c, http.StatusOK, result)
	// } else {
	// 	var errMsg string
	// 	switch err.Error() {
	// 	case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
	// 		errMsg = err.Error()
	// 	default:
	// 		errMsg = utils.BadRequest
	// 	}
	// 	result := utils.NewHTTPError(errMsg)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}

func (handler *CollectionHandler) DeleteComment(c *gin.Context) {
	// var requestBody model.Comment
	// c.BindJSON(&requestBody)

	// if !requestBody.UserID.IsZero() {
	// 	err := handler.comment.DeleteComment(c.Request.Context(), utils.Session, c.Param("commentId"), requestBody.UserID.Hex())
	// 	if err == nil {
	// 		result := utils.SuccessPayload(nil, "Successfully deleted comment")
	// 		utils.Response(c, http.StatusOK, result)
	// 	} else {
	// 		var errMsg string
	// 		switch err.Error() {
	// 		case utils.CommentDoesNotExist, utils.CommentCannotUpdateOtherUserComment:
	// 			errMsg = err.Error()
	// 		default:
	// 			errMsg = utils.BadRequest
	// 		}
	// 		result := utils.NewHTTPError(errMsg)
	// 		utils.Response(c, http.StatusBadRequest, result)
	// 	}
	// } else {
	// 	result := utils.NewHTTPError(utils.UserIdMissing)
	// 	utils.Response(c, http.StatusBadRequest, result)
	// }
}
