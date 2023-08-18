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

// CreatePost creates new post
//
//	@Summary		create new post
//	@Description	create new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.CreatePostRequest	true	"Create post parameters"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/ [post]
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
	resp, err := svc.authenticateAndPostClient.CreatePost(ctx, &pb_aap.CreatePostRequest{
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

// GetPostDetailInfo gets post detail information
//
//	@Summary		get post detail information
//	@Description	get post detail information
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	types.PostDetailInfoResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/{post_id} [get]
func (svc *WebService) GetPostDetailInfo(ctx *gin.Context) {
	// Check URL params
	stringPostId := ctx.Param("post_id")
	postId, err := strconv.Atoi(stringPostId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call gprc service
	resp, err := svc.authenticateAndPostClient.GetPostDetailInfo(ctx, &pb_aap.GetPostDetailInfoRequest{
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
		ctx.IndentedJSON(http.StatusAccepted, svc.newPostDetailInfoResponse(resp.Post))
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

// EditPost edits post information
//
//	@Summary		edit post information
//	@Description	edit post information
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/ [put]
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

	resp, err := svc.authenticateAndPostClient.EditPost(ctx, &pb_aap.EditPostRequest{
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

// DeletePost deletes post
//
//	@Summary		delete post
//	@Description	delete post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/ [delete]
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
	resp, err := svc.authenticateAndPostClient.DeletePost(ctx, &pb_aap.DeletePostRequest{
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

// CommentPost comments to a post
//
//	@Summary		comment to post
//	@Description	comment to post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int								true	"Post ID"
//	@Param			request	body		types.CreatePostCommentRequest	true	"Comment's content"
//	@Success		200		{object}	types.PostDetailInfoResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/{post_id} [post]
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
	resp, err := svc.authenticateAndPostClient.CommentPost(ctx,
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

// LikePost likes a post
//
//	@Summary		like post
//	@Description	like post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/posts/{post_id}/likes [post]
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
	resp, err := svc.authenticateAndPostClient.LikePost(ctx,
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
	} else if resp.GetStatus() == pb_aap.LikePostResponse_ALREADY_LIKED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "already liked"})
		return	
	} else if resp.GetStatus() == pb_aap.LikePostResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) newPostDetailInfoResponse(post *pb_aap.PostDetailInfo) types.PostDetailInfoResponse {
	var comments []types.CommentResponse
	for _, comment := range post.GetComments() {
		comments = append(comments, types.CommentResponse{
			CommentId:   comment.GetCommentId(),
			UserId:      comment.GetUserId(),
			PostId:      comment.GetPostId(),
			ContentText: comment.GetContentText(),
		})
	}

	var users_liked []int64
	for _, like := range post.GetLikedUsers() {
		users_liked = append(users_liked, like.GetUserId())
	}

	return types.PostDetailInfoResponse{
		PostID:           post.GetPostId(),
		UserID:           post.GetUserId(),
		ContentText:      post.GetContentText(),
		ContentImagePath: post.GetContentImagePath(),
		CreatedAt:        post.GetCreatedAt().AsTime().In(time.Local).Format(time.DateTime),
		Comments:         comments,
		UsersLiked:       users_liked,
	}
}
