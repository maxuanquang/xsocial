package v1

import (
	"github.com/gin-gonic/gin"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

var (
	aapClient pb_aap.AuthenticateAndPostClient
	// nfClient  pb_nf.NewsfeedClient
)

func AddAllRouter(r *gin.RouterGroup, in_aapClient pb_aap.AuthenticateAndPostClient, in_nfClient pb_nf.NewsfeedClient) {
	aapClient = in_aapClient
	// nfClient = in_nfClient

	AddUserRouter(r)
}
