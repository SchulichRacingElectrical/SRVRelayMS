package handlers

import (
	"context"
	"database-ms/app/middleware"
	"database-ms/app/model"
	"database-ms/app/services"
	utils "database-ms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
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

func (handler *CollectionHandler) AddComment(ctx *gin.Context) {
	var newComment model.CollectionComment
	err := ctx.BindJSON(&newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-tenant write
	organization, _ := middleware.GetOrganizationClaim(ctx)
	collection, perr := handler.collectionService.FindById(ctx, newComment.CollectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	newComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the collection
	err = handler.collectionService.AddComment(ctx.Request.Context(), &newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
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

	// Guard against cross-tenant read
	organization, _ := middleware.GetOrganizationClaim(ctx)
	collection, perr := handler.collectionService.FindById(ctx, collectionId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
		return
	}
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.ThingNotFound))
		return
	}
	if thing.OrganizationId != organization.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to read the comments
	comments, err := handler.collectionService.GetComments(ctx.Request.Context(), collectionId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) UpdateCommentContent(ctx *gin.Context) {
	var updatedComment model.CollectionComment
	err := ctx.BindJSON(&updatedComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Guard against cross-user update
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}
	updatedComment.UserId = user.Id
	comment, err := handler.collectionService.GetComment(ctx, updatedComment.Id)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	updatedComment.LastUpdate = utils.CurrentTimeInMilli()

	// Attempt to create the collection
	err = handler.collectionService.UpdateCommentContent(ctx.Request.Context(), &updatedComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
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

	// Guard against cross-user update
	user, err := middleware.GetUserClaim(ctx)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}
	comment, err := handler.collectionService.GetComment(ctx, commentId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Attempt to delete the comment
	err = handler.collectionService.DeleteComment(ctx.Request.Context(), commentId)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CollectionHandler) IsUserPartofOrg(collectionId uuid.UUID, userOrgId uuid.UUID, ctx context.Context) (bool, *pgconn.PgError) {
	collection, perr := handler.collectionService.FindById(ctx, collectionId)
	if perr != nil {
		return false, perr
	}
	thing, perr := handler.thingService.FindById(ctx, collection.ThingId)
	if perr != nil {
		return false, perr
	}

	// Guard against cross-tenant deletion
	if thing.OrganizationId != userOrgId {
		return false, nil
	}

	return true, nil
}
