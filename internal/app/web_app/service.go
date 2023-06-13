package web_app

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/maxuanquang/social-network/internal/app/web_app/v1"
)

func Run() {
	r := gin.Default()
	v1Router := r.Group("v1")
	v1.AddUserRouter(v1Router)
	r.Run(":8080")
}