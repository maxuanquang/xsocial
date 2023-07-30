package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

// AddNewsfeedRouter adds newsfeed-related routes to input router
func AddNewsfeedRouter(r *gin.RouterGroup, svc *service.WebService) {
	postRouter := r.Group("newsfeed")

	postRouter.GET("", svc.GetNewsfeed)
}
