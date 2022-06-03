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

type CommentHandler struct {
	commentService    services.CommentServiceInterface
	thingService      services.ThingServiceInterface
	sessionService    services.SessionServiceInterface
	collectionService services.CollectionServiceInterface
}

func NewCommentAPI(
	commentService services.CommentServiceInterface,
	thingService services.ThingServiceInterface,
	sessionService services.SessionServiceInterface,
	collectionService services.CollectionServiceInterface,
) *CommentHandler {
	return &CommentHandler{
		commentService:    commentService,
		thingService:      thingService,
		sessionService:    sessionService,
		collectionService: collectionService,
	}
}

func (handler *CommentHandler) CreateComment(ctx *gin.Context) {
	// Attempt to parse the body
	var newComment model.Comment
	err := ctx.BindJSON(&newComment)
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.BadRequest))
		return
	}

	// Attempt to find the associated thingId
	thingId := uuid.UUID{}
	if newComment.CollectionId != nil {
		collection, perr := handler.collectionService.FindById(ctx, *newComment.CollectionId)
		if perr != nil {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CollectionNotFound))
			return
		}
		thingId = collection.ThingId
	} else {
		session, perr := handler.sessionService.FindById(ctx, *newComment.SessionId)
		if perr != nil {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.SessionNotFound))
			return
		}
		thingId = session.ThingId
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, thingId)
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

	// Update the comment's time
	newComment.Time = utils.CurrentTimeInMilli()

	// Attempt to create the comment
	perr = handler.commentService.CreateComment(ctx.Request.Context(), &newComment)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(newComment, "Successfully created comment.")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CommentHandler) GetComments(ctx *gin.Context) {
	// Attempt to read from the params
	contextId, err := uuid.Parse(ctx.Param("contextId"))
	if err != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Attempt to find the associated thingId
	thingId := uuid.UUID{}
	collection, perr := handler.collectionService.FindById(ctx, contextId)
	if perr != nil {
		session, perr := handler.sessionService.FindById(ctx, contextId)
		if perr != nil {
			utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError("")) // TODO: Comment context not found
			return
		}
		thingId = session.ThingId
	} else {
		thingId = collection.ThingId
	}

	// Attempt to find the thing
	thing, perr := handler.thingService.FindById(ctx, thingId)
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
	comments, perr := handler.commentService.FindCommentsByContextId(ctx.Request.Context(), contextId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentsNotFound))
		return
	}

	// Send the response
	result := utils.SuccessPayload(comments, "Successfully retrieved comments")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CommentHandler) UpdateComment(ctx *gin.Context) {
	// Attempt to parse the body
	var updatedComment model.Comment
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
	comment, perr := handler.commentService.FindById(ctx, updatedComment.Id)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPError(utils.CommentNotFound))
		return
	}

	// Guard against cross-tenant updates
	if comment.UserId != user.Id {
		utils.Response(ctx, http.StatusUnauthorized, utils.NewHTTPError(utils.Unauthorized))
		return
	}

	// Update only the comments and remove the sub comments
	comment.Content = updatedComment.Content
	comment.Comments = []model.Comment{}

	// Attempt to update the comment
	perr = handler.commentService.UpdateComment(ctx.Request.Context(), comment)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, perr.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully updated")
	utils.Response(ctx, http.StatusOK, result)
}

func (handler *CommentHandler) DeleteComment(ctx *gin.Context) {
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
	comment, perr := handler.commentService.FindById(ctx, commentId)
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
	perr = handler.commentService.DeleteComment(ctx.Request.Context(), commentId)
	if perr != nil {
		utils.Response(ctx, http.StatusBadRequest, utils.NewHTTPCustomError(utils.BadRequest, err.Error()))
		return
	}

	// Send the response
	result := utils.SuccessPayload(nil, "Successfully deleted")
	utils.Response(ctx, http.StatusOK, result)
}
