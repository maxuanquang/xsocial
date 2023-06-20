package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"

	// client_nf "github.com/maxuanquang/social-network/pkg/client/newsfeed"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

var validate = types.NewValidator()

type WebService struct {
	AuthenticateAndPostClient pb_aap.AuthenticateAndPostClient
	NewsfeedClient            pb_nf.NewsfeedClient
	RedisClient               *redis.Client
}

func NewWebService(conf *configs.WebConfig) (*WebService, error) {
	aapClient, err := client_aap.NewClient(conf.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// nfClient, err := client_nf.NewClient(conf.Newsfeed.Hosts)
	// if err != nil {
	// 	return nil, err
	// }

	redisClient := redis.NewClient(&redis.Options{Addr: conf.Redis.Addr, Password: conf.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	return &WebService{
		AuthenticateAndPostClient: aapClient,
		// NewsfeedClient:            nfClient,
		RedisClient: redisClient,
	}, nil
}

// CheckUserAuthentication checks authentication of user and provides
// a session_id cookie if authentication succeeded
func (svc *WebService) CheckUserAuthentication(ctx *gin.Context) {
	var jsonRequest types.LoginRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("jsonRequest binds error: %v", err.Error())})
		return
	}

	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authentication, err := svc.AuthenticateAndPostClient.CheckUserAuthentication(ctx, &pb_aap.UserInfo{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// If logged in, set a sessionId for this session
	sessionId := uuid.New().String()

	// Save current sessionID and username in Redis
	err = svc.RedisClient.Set(svc.RedisClient.Context(), sessionId, authentication.GetInfo().GetUserName(), 300*time.Second).Err()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set sessionID cookie
	ctx.SetCookie("session_id", sessionId, 300, "", "", false, false)

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Log in succeeded"})
}

func (svc *WebService) CreateUser(ctx *gin.Context) {
	var jsonRequest types.CreateUserRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Call CreateUser service
	dob, _ := time.Parse(time.DateOnly, jsonRequest.DoB)
	_, err = svc.AuthenticateAndPostClient.CreateUser(ctx, &pb_aap.UserDetailInfo{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
		FirstName:    jsonRequest.FirstName,
		LastName:     jsonRequest.LastName,
		Dob:          dob.Unix(),
		Email:        jsonRequest.Email,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Create user failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Create user succeeded"})
}

func (svc *WebService) EditUser(ctx *gin.Context) {
	// Check authorization
	sessionId, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	userName, err := svc.RedisClient.Get(svc.RedisClient.Context(), sessionId).Result()
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Check request format
	var jsonRequest types.EditUserRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Call EditUser service
	dob, _ := time.Parse(time.DateOnly, jsonRequest.DoB)
	_, err = svc.AuthenticateAndPostClient.EditUser(ctx, &pb_aap.UserDetailInfo{
		UserName:     userName,
		UserPassword: jsonRequest.Password,
		FirstName:    jsonRequest.FirstName,
		LastName:     jsonRequest.LastName,
		Dob:          dob.Unix(),
		Email:        jsonRequest.Email,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Edit user failed: %v", err)})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Edit user succeeded"})

}
