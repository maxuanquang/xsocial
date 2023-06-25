package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (svc *WebService) GetUserFollower(ctx *gin.Context) {
	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call GetUserFollower gprc service
	userFollower, err := svc.AuthenticateAndPostClient.GetUserFollower(ctx, &pb_aap.UserInfo{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Get user follower failed: %v", err)})
		return
	}

	// Return necessary information
	var followers []map[string]interface{}
	for _, follower := range userFollower.GetFollowers() {
		followers = append(followers, map[string]interface{}{"id": follower.UserId, "username": follower.UserName})
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Get followers succeeded", "followers": followers})
}

func (svc *WebService) FollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call FollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.FollowUser(ctx,
		&pb_aap.UserAndFollowerInfo{
			User:     &pb_aap.UserInfo{UserId: int64(userId)},
			Follower: &pb_aap.UserInfo{UserId: int64(followerId)},
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call UnfollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.UnfollowUser(ctx,
		&pb_aap.UserAndFollowerInfo{
			User:     &pb_aap.UserInfo{UserId: int64(userId)},
			Follower: &pb_aap.UserInfo{UserId: int64(followerId)},
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) GetUserPost(ctx *gin.Context) {
	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call GetUserPost grpc service
	userPosts, err := svc.AuthenticateAndPostClient.GetUserPost(ctx,
		&pb_aap.UserInfo{
			UserId: int64(userId),
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Return
	var posts []gin.H
	for _, postDetailInfo := range userPosts.Posts {
		posts = append(posts, svc.newJSONPost(postDetailInfo))
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK", "posts": posts})
}
