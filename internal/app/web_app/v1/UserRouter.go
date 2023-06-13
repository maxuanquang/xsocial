package v1

import "github.com/gin-gonic/gin"

// AddUserRouter adds user-related routes to input router
func AddUserRouter(r *gin.RouterGroup) {
	userRouter := r.Group("users")
	userRouter.GET("", func(ctx *gin.Context) {})
	userRouter.POST("", func(ctx *gin.Context) {})
	userRouter.PUT("", func(ctx *gin.Context) {})
}