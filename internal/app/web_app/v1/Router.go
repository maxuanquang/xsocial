package v1

import (
	"github.com/gin-gonic/gin"
	// pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	// pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

func AddAllRouter(r *gin.RouterGroup, svc *service.WebService) {
	AddUserRouter(r, svc)
	AddFriendRouter(r, svc)
	AddPostRouter(r, svc)
	AddNewsfeedRouter(r, svc)
}
