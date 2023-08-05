package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (svc *WebService) GetUserFollower(ctx *gin.Context) {
	// Validate parameter
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserFollower gprc service
	resp, err := svc.AuthenticateAndPostClient.GetUserFollower(ctx, &pb_aap.GetUserFollowerRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return
	ctx.IndentedJSON(http.StatusAccepted, resp.GetFollowersIds())
}

func (svc *WebService) GetUserFollowing(ctx *gin.Context) {
	// Validate parameter
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserFollower gprc service
	resp, err := svc.AuthenticateAndPostClient.GetUserFollowing(ctx, &pb_aap.GetUserFollowingRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return
	ctx.IndentedJSON(http.StatusAccepted, resp.GetFollowingsIds())
}

func (svc *WebService) FollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: "session unauthorized"})
		return
	}

	// Validate parameter
	followingId, err := strconv.Atoi(ctx.Param("following_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call FollowUser grpc service
	resp, err := svc.AuthenticateAndPostClient.FollowUser(ctx,
		&pb_aap.FollowUserRequest{
			UserId:      int64(userId),
			FollowingId: int64(followingId),
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.FollowUserResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.FollowUserResponse_ALREADY_FOLLOWED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "already followed"})
		return
	} else if resp.GetStatus() == pb_aap.FollowUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unkown error"})
		return
	}
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: "session unauthorized"})
		return
	}

	// Validate parameter
	follwingId, err := strconv.Atoi(ctx.Param("following_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call UnfollowUser grpc service
	resp, err := svc.AuthenticateAndPostClient.UnfollowUser(ctx,
		&pb_aap.UnfollowUserRequest{
			UserId:      int64(userId),
			FollowingId: int64(follwingId)},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.UnfollowUserResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.UnfollowUserResponse_NOT_FOLLOWED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "not followed"})
		return
	} else if resp.GetStatus() == pb_aap.UnfollowUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unkown error"})
		return
	}
}

func (svc *WebService) GetUserPosts(ctx *gin.Context) {
	// Validate parameter
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	}

	// Call GetUserPost grpc service
	resp, err := svc.AuthenticateAndPostClient.GetUserPosts(ctx,
		&pb_aap.GetUserPostsRequest{
			UserId: int64(userId),
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserPostsResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
	} else if resp.GetStatus() == pb_aap.GetUserPostsResponse_OK {
		ctx.IndentedJSON(http.StatusOK, resp.GetPostsIds())
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}

}
