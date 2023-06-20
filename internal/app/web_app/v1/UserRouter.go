package v1

import (
	// "context"
	// "fmt"
	// "net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
	// pb "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

// AddUserRouter adds user-related routes to input router
func AddUserRouter(r *gin.RouterGroup, svc *service.WebService) {
	userRouter := r.Group("users")
	userRouter.POST("register", svc.CreateUser)
	userRouter.POST("login", svc.CheckUserAuthentication)
	userRouter.POST("edit", svc.EditUser)
}

// func editHandler(ctx *gin.Context) {
// 	username := ctx.PostForm("username")
// 	password := ctx.PostForm("password")
// 	firstname := ctx.PostForm("firstname")
// 	lastname := ctx.PostForm("lastname")
// 	dob := ctx.PostForm("dob")
// 	email := ctx.PostForm("email")

// 	// Process the contents
// 	var dob_unix int64

// 	if dob != "" {
// 		processedDOB, err := time.Parse(time.DateOnly, dob)
// 		if err != nil {
// 			ctx.IndentedJSON(
// 				http.StatusForbidden,
// 				gin.H{
// 					"message": "dob needs to have following format yyyy-mm-dd",
// 					"error":   fmt.Sprintf("edit user information failed: %v", err),
// 				},
// 			)
// 			return
// 		}
// 		dob_unix = processedDOB.Unix()
// 	}

// 	_, err := aapClient.EditUser(
// 		context.Background(),
// 		&pb.UserDetailInfo{
// 			UserName:     username,
// 			UserPassword: password,
// 			FirstName:    firstname,
// 			LastName:     lastname,
// 			Dob:          dob_unix,
// 			Email:        email,
// 		},
// 	)
// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("edit user information failed: %v", err)})
// 		return
// 	}

// 	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "edit user information succeeded!"})
// }

// func createUserHandler(ctx *gin.Context) {
// 	username := ctx.PostForm("username")
// 	password := ctx.PostForm("password")
// 	firstname := ctx.PostForm("firstname")
// 	lastname := ctx.PostForm("lastname")
// 	dob := ctx.PostForm("dob")
// 	email := ctx.PostForm("email")

// 	// Process the contents
// 	processedDOB, err := time.Parse(time.DateOnly, dob)
// 	if err != nil {
// 		ctx.IndentedJSON(
// 			http.StatusForbidden,
// 			gin.H{
// 				"message": "dob needs to have following format yyyy-mm-dd",
// 				"error":   fmt.Sprintf("create user failed: %v", err),
// 			},
// 		)
// 		return
// 	}

// 	_, err = aapClient.CreateUser(
// 		context.Background(),
// 		&pb.UserDetailInfo{
// 			UserName:     username,
// 			UserPassword: password,
// 			FirstName:    firstname,
// 			LastName:     lastname,
// 			Dob:          processedDOB.Unix(),
// 			Email:        email,
// 		},
// 	)
// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("create user failed: %v", err)})
// 		return
// 	}

// 	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "create user succeeded!"})
// }

// func loginHandler(ctx *gin.Context) {
// 	username := ctx.PostForm("username")
// 	password := ctx.PostForm("password")

// 	_, err := aapClient.CheckUserAuthentication(context.Background(), &pb.UserInfo{UserId: 1, UserName: username, UserPassword: password})
// 	if err != nil {
// 		ctx.IndentedJSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Login failed: %v", err)})
// 		return
// 	}

// 	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Log in succeeded!"})
// }
