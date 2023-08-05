package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

// AddUserRouter adds user-related routes to input router
func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("users")
	userRouter.POST("signup", svc.CreateUser)
	userRouter.POST("login", svc.CheckUserAuthentication)
	userRouter.POST("edit", svc.EditUser)
	userRouter.GET(":user_id", svc.GetUserDetailInfo)
}
