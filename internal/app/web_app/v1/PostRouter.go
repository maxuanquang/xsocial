package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

// AddPostRouter adds post-related routes to input router
func AddPostRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("posts")

	postRouter.POST("", svc.CreatePost)
	postRouter.GET("/url", svc.GetS3PresignedUrl)
	postRouter.GET(":post_id", svc.GetPostDetailInfo)
	postRouter.PUT(":post_id", svc.EditPost)
	postRouter.DELETE(":post_id", svc.DeletePost)

	postRouter.POST(":post_id/comments", svc.CommentPost)
	postRouter.POST(":post_id/likes", svc.LikePost)
}
