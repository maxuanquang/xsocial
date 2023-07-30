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
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user_id does not exist"})
		return
	}

	// Call GetUserFollower gprc service
	userFollower, err := svc.AuthenticateAndPostClient.GetUserFollower(ctx, &pb_aap.UserInfo{
		Id: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return necessary information
	var followers []map[string]interface{}
	for _, follower := range userFollower.GetFollowers() {
		followers = append(followers, map[string]interface{}{"id": follower.GetId(), "username": follower.GetUserName()})
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"followers": followers})
}

func (svc *WebService) FollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user_id does not exist"})
		return
	}

	// Call FollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.FollowUser(ctx,
		&pb_aap.UserAndFollowerInfo{
			User:     &pb_aap.UserInfo{Id: int64(userId)},
			Follower: &pb_aap.UserInfo{Id: int64(followerId)},
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user_id does not exist"})
		return
	}

	// Call UnfollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.UnfollowUser(ctx,
		&pb_aap.UserAndFollowerInfo{
			User:     &pb_aap.UserInfo{Id: int64(userId)},
			Follower: &pb_aap.UserInfo{Id: int64(followerId)},
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) GetUserPost(ctx *gin.Context) {
	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user_id does not exist"})
		return
	}

	// Call GetUserPost grpc service
	userPosts, err := svc.AuthenticateAndPostClient.GetUserPost(ctx,
		&pb_aap.UserInfo{
			Id: int64(userId),
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return
	var posts []gin.H
	for _, postDetailInfo := range userPosts.Posts {
		posts = append(posts, svc.newMapPost(postDetailInfo))
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"posts": posts})
}
