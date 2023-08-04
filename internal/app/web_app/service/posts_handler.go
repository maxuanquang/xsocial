package service

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/pkg/types"

	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (svc *WebService) CreatePost(ctx *gin.Context) {
	// Check session
	_, userId, err := svc.checkSessionAuthentication(ctx)
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
	var visible = true
	if jsonRequest.Visible != nil && !*jsonRequest.Visible {
		visible = false
	}
	resp, err := svc.AuthenticateAndPostClient.CreatePost(ctx, &pb_aap.CreatePostRequest{
		UserId:           int64(userId),
		ContentText:      jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible:          visible,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CreatePostResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.CreatePostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) GetPostDetailInfo(ctx *gin.Context) {
	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call gprc service
	resp, err := svc.AuthenticateAndPostClient.GetPostDetailInfo(ctx, &pb_aap.GetPostDetailInfoRequest{
		PostId: int64(postId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetPostDetailInfoResponse_POST_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetPostDetailInfoResponse_OK {
		ctx.IndentedJSON(http.StatusAccepted, svc.newMapPost(resp.Post))
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) EditPost(ctx *gin.Context) {
	// Check session
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
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
	var contentText *string
	if jsonRequest.ContentText != nil {
		contentText = jsonRequest.ContentText
	}
	var contentImagePath *string
	if jsonRequest.ContentImagePath != nil {
		stringLinks := strings.Join(*jsonRequest.ContentImagePath, " ")
		contentImagePath = &stringLinks
	}
	var visible *bool
	if jsonRequest.Visible != nil {
		visible = jsonRequest.Visible
	}

	resp, err := svc.AuthenticateAndPostClient.EditPost(ctx, &pb_aap.EditPostRequest{
		UserId:           int64(userId),
		PostId:           int64(postId),
		ContentText:      contentText,
		ContentImagePath: contentImagePath,
		Visible:          visible,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.EditPostResponse_POST_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_NOT_ALLOWED {
		ctx.IndentedJSON(http.StatusForbidden, types.MessageResponse{Message: "not allowed"})
		return
	} else if resp.GetStatus() == pb_aap.EditPostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) DeletePost(ctx *gin.Context) {
	// Check session
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call grpc service
	resp, err := svc.AuthenticateAndPostClient.DeletePost(ctx, &pb_aap.DeletePostRequest{
		PostId: int64(postId),
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.DeletePostResponse_POST_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_NOT_ALLOWED {
		ctx.IndentedJSON(http.StatusForbidden, types.MessageResponse{Message: "not allowed"})
		return
	} else if resp.GetStatus() == pb_aap.DeletePostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) CommentPost(ctx *gin.Context) {
	// Check session
	_, userId, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Check URL params
	postId, err := strconv.Atoi(ctx.Param("post_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Check request
	var jsonRequest types.CreatePostCommentRequest
	err = ctx.BindJSON(&jsonRequest)
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
	resp, err := svc.AuthenticateAndPostClient.CommentPost(ctx,
		&pb_aap.CommentPostRequest{
			PostId:      int64(postId),
			UserId:      int64(userId),
			ContentText: jsonRequest.ContentText,
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CommentPostResponse_POST_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.CommentPostResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return

	} else if resp.GetStatus() == pb_aap.CommentPostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) LikePost(ctx *gin.Context) {
	// Check session
	_, userId, err := svc.checkSessionAuthentication(ctx)
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
	resp, err := svc.AuthenticateAndPostClient.LikePost(ctx,
		&pb_aap.LikePostRequest{
			PostId: int64(postId),
			UserId: int64(userId),
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.LikePostResponse_POST_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	} else if resp.GetStatus() == pb_aap.LikePostResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return

	} else if resp.GetStatus() == pb_aap.LikePostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) newMapPost(post *pb_aap.PostDetailInfo) gin.H {
	var comments []map[string]interface{}
	for _, comment := range post.GetComments() {
		comments = append(comments, map[string]interface{}{
			"comment_id":   comment.GetCommentId(),
			"user_id":      comment.GetUserId(),
			"post_id":      comment.GetPostId(),
			"content_text": comment.GetContentText(),
		})
	}

	var likes []map[string]interface{}
	for _, like := range post.GetLikedUsers() {
		likes = append(likes, map[string]interface{}{
			"user_id": like.GetUserId(),
			"post_id": like.GetPostId(),
		})
	}

	return gin.H{
		"post_id":            post.GetPostId(),
		"created_at":         post.GetCreatedAt().AsTime().In(time.Local).Format(time.DateTime),
		"user_id":            post.GetUserId(),
		"content_text":       post.GetContentText(),
		"content_image_path": post.GetContentImagePath(),
		"comments":           comments,
		"likes":              likes,
	}
}
