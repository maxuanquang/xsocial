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
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.CreatePostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) GetPost(ctx *gin.Context) {
	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong post_id"})
		return
	}

	// Call gprc service
	postDetailInfo, err := svc.AuthenticateAndPostClient.GetPost(ctx, &pb_aap.PostInfo{PostId: int64(postId)})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something is wrong"})
		return
	}
	if !postDetailInfo.GetVisible() {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Post is deleted"})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, svc.newJSONPost(postDetailInfo))
}

func (svc *WebService) EditPost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong post_id"})
		return
	}

	// Check EditPostRequest
	var jsonRequest types.EditPostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.EditPost(ctx, &pb_aap.PostDetailInfo{
		UserId:           int64(userId),
		PostId:           int64(postId),
		ContentText:      jsonRequest.ContentText,
		ContentImagePath: jsonRequest.ContentImagePath,
		Visible:          true,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) DeletePost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong post_id"})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.DeletePost(ctx, &pb_aap.PostInfo{
		PostId: int64(postId),
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

// POST: ver/posts/:post_id/comment -> msg
func (svc *WebService) CommentPost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong post_id"})
		return
	}

	// Check request
	var jsonRequest types.CommentPostRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) LikePost(ctx *gin.Context) {
	// Check session
	_, userId, _, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong post_id"})
		return
	}

	// Call grpc service
	_, err = svc.AuthenticateAndPostClient.LikePost(ctx,
		&pb_aap.LikeInfo{
			PostId: int64(postId),
			UserId: int64(userId),
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return message
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) newJSONPost(postDetailInfo *pb_aap.PostDetailInfo) gin.H {
	var comments []map[string]interface{}
	for _, comment := range postDetailInfo.GetComments() {
		comments = append(comments, map[string]interface{}{
			"commentId": comment.GetCommentId(),
			"userId":    comment.GetUserId(),
			"postId":    comment.GetPostId(),
			"content":   comment.GetContent(),
		})
	}

	var likes []map[string]interface{}
	for _, like := range postDetailInfo.GetLikes() {
		likes = append(likes, map[string]interface{}{
			"userId": like.GetUserId(),
			"postId": like.GetPostId(),
		})
	}

	return gin.H{
		"postId":             postDetailInfo.GetPostId(),
		"userId":             postDetailInfo.GetUserId(),
		"content_text":       postDetailInfo.GetContentText(),
		"content_image_path": postDetailInfo.GetContentImagePath(),
		"create_at":          postDetailInfo.GetCreatedAt(),
		"comments":           comments,
		"likes":              likes,
	}
}
