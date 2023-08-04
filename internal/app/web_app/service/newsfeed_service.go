package service

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/maxuanquang/social-network/internal/pkg/types"
// 	"net/http"

// 	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
// )

// func (svc *WebService) GetNewsfeed(ctx *gin.Context) {
// 	// Check authorization
// 	_, userId, _, err := svc.checkSessionAuthentication(ctx)
// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
// 		return
// 	}

// 	// Call GetNewsfeed service
// 	newsfeed, err := svc.NewsfeedClient.GetNewsfeed(ctx, &pb_nf.NewsfeedRequest{UserId: int64(userId)})
// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
// 		return
// 	}

// 	// Return
// 	var posts []gin.H
// 	for _, postDetailInfo := range newsfeed.Posts {
// 		posts = append(posts, svc.newMapPost(postDetailInfo))
// 	}

// 	ctx.IndentedJSON(http.StatusOK, gin.H{"newsfeed": posts})
// }
