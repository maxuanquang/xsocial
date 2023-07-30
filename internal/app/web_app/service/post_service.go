package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/pkg/types"

	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (svc *WebService) CreatePost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.CreatePostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.CreatePost(ctx, &pb_aap.PostDetailInfo{
		UserId:           int64(userId),
		ContentText:      jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible:          true,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) GetPost(ctx *gin.Context) {
	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "Wrong post_id"})
		return
	}

	// Call gprc service
	postDetailInfo, err := svc.AuthenticateAndPostClient.GetPost(ctx, &pb_aap.PostInfo{Id: int64(postId)})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "Something is wrong"})
		return
	}
	if !postDetailInfo.GetVisible() {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "Post is not available"})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, svc.newMapPost(postDetailInfo))
}

func (svc *WebService) EditPost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "Wrong post_id"})
		return
	}

	// Check EditPostRequest
	var jsonRequest types.EditPostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.EditPost(ctx, &pb_aap.PostDetailInfo{
		Id:               int64(postId),
		UserId:           int64(userId),
		ContentText:      jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible:          true,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) DeletePost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "Wrong post_id"})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.DeletePost(ctx, &pb_aap.PostInfo{
		Id:     int64(postId),
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

// POST: ver/posts/:post_id/comment -> msg
func (svc *WebService) CommentPost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "Wrong post_id"})
		return
	}

	// Check request
	var jsonRequest types.CommentPostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.CommentPost(ctx,
		&pb_aap.CommentInfo{
			PostId:  int64(postId),
			UserId:  int64(userId),
			Content: jsonRequest.Content,
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) LikePost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "Wrong post_id"})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.LikePost(ctx,
		&pb_aap.LikeInfo{
			PostId: int64(postId),
			UserId: int64(userId),
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) newMapPost(postDetailInfo *pb_aap.PostDetailInfo) gin.H {
	var comments []map[string]interface{}
	for _, comment := range postDetailInfo.GetComments() {
		comments = append(comments, map[string]interface{}{
			"id":      comment.GetId(),
			"user_id": comment.GetUserId(),
			"post_id": comment.GetPostId(),
			"content": comment.GetContent(),
		})
	}

	var likes []map[string]interface{}
	for _, like := range postDetailInfo.GetLikes() {
		likes = append(likes, map[string]interface{}{
			"user_id": like.GetUserId(),
			"post_id": like.GetPostId(),
		})
	}

	return gin.H{
		"id":                 postDetailInfo.GetId(),
		"user_id":            postDetailInfo.GetUserId(),
		"content_text":       postDetailInfo.GetContentText(),
		"content_image_path": postDetailInfo.GetContentImagePath(),
		"create_at":          postDetailInfo.GetCreatedAt(),
		"comments":           comments,
		"likes":              likes,
	}
}
