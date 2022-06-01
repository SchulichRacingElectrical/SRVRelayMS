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
	// Attempt to parse the body
	var newCollection model.Collection
	err := ctx.BindJSON(&newCollection)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against non-lead+
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to find the associated thing
	thing, perr := handler.thingService.FindById(ctx, newCollection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant writes
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if organization.Id != thing.OrganizationId {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to create the collection
	perr = handler.collectionService.CreateCollection(ctx.Request.Context(), &newCollection)
	if perr != nil {
		if perr.Code == "23505" {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError("")) // TODO: Add conflict error
		} else {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
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
	thing, perr := handler.thingService.FindById(ctx.Request.Context(), thingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reading
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the collections
	collections, perr := handler.collectionService.FindCollectionsByThingId(ctx.Request.Context(), thingId)
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

	// Guard against non-lead+ requests
	if !middleware.IsAuthorizationAtLeast(ctx, "Lead") {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to get the thing
	organization, _ := middleware.GetOrganizationClaim(ctx)
	thing, perr := handler.thingService.FindById(ctx, updatedCollection.ThingId)
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
			utils.Response(ctx, http.StatusConflict, utils.NewHTTPError("")) // TODO: Not unique error
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
	collection, perr := handler.collectionService.FindById(ctx.Request.Context(), collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant deletion
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the collection
	perr = handler.collectionService.DeleteCollection(ctx.Request.Context(), collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) AddComment(ctx *gin.Context) {
	// Attempt to parse the body
	var newComment model.CollectionComment
	err := ctx.BindJSON(&newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the collection
	collection, perr := handler.collectionService.FindById(ctx, newComment.CollectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant writes
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Update the comment's update time
	newComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the comment
	perr = handler.collectionService.CreateComment(ctx.Request.Context(), &newComment)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newComment, "Successfully added comment.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) GetComments(ctx *gin.Context) {
	// Attempt to read from the params
	collectionId, err := uuid.Parse(ctx.Param("collectionId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the collection
	collection, perr := handler.collectionService.FindById(ctx, collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}

	// Guard against cross-tenant reads
	organization, _ := middleware.GetOrganizationClaim(ctx)
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the comments
	comments, perr := handler.collectionService.FindCommentsByCollectionId(ctx.Request.Context(), collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) UpdateCommentContent(ctx *gin.Context) {
	// Attempt to parse the body
	var updatedComment model.CollectionComment
	err := ctx.BindJSON(&updatedComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the user
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the comment
	comment, perr := handler.collectionService.FindCommentById(ctx, updatedComment.Id)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}

	// Guard against cross-tenant updates
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Update the last update of the comment
	comment.LastUpdate = utils.CurrentTimeInMilli()
	comment.Content = updatedComment.Content

	// Attempt to update the comment
	perr = handler.collectionService.UpdateCommentContent(ctx.Request.Context(), comment)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) DeleteComment(ctx *gin.Context) {
	// Attempt to read from the params
	commentId, err := uuid.Parse(ctx.Param("commentId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the user
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the comment
	comment, perr := handler.collectionService.FindCommentById(ctx, commentId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}

	// Guard against cross-tenant deletions
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the comment
	perr = handler.collectionService.DeleteComment(ctx.Request.Context(), commentId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}
