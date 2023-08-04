package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

func AddAllRouter(r *gin.RouterGroup, svc *service.WebService) {
	AddUserRouter(r, svc)
	AddFriendRouter(r, svc)
	AddPostRouter(r, svc)
	// AddNewsfeedRouter(r, svc)
}
