package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/pkg/types"

	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

// GetNewsfeed gets user's newsfeed
//
//	@Summary		get user's newsfeed
//	@Description	get user's newsfeed
//	@Tags			newsfeed
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.NewsfeedResponse
//	@Failure		400	{object}	types.MessageResponse
//	@Failure		500	{object}	types.MessageResponse
//	@Router			/newsfeed [get]
func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call GetNewsfeed service
	resp, err := svc.newsfeedClient.GetNewsfeed(ctx, &pb_nf.GetNewsfeedRequest{UserId: int64(userId)})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "newsfeed empty"})
		return
	} else if resp.GetStatus() == pb_nf.GetNewsfeedResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.NewsfeedResponse{PostsIds: resp.GetPostsIds()})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}
