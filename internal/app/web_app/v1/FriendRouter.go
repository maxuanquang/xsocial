package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

// AddFriendRouter adds friend-related routes to input router
func AddFriendRouter(r *gin.RouterGroup, svc *service.WebService) {
	friendRouter := r.Group("friends")
	friendRouter.GET(":user_id", svc.GetUserFollower)
	friendRouter.POST(":user_id", svc.FollowUser)
	friendRouter.DELETE(":user_id", svc.UnfollowUser)
}
