package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

func (svc *WebService) CheckUserAuthentication(ctx *gin.Context) {
	// Validate request
	var jsonRequest types.LoginRequest
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

	// Call CheckUserAuthentication service
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

	// Save current sessionID and expiration time in Redis
	err = svc.RedisClient.HSet(svc.RedisClient.Context(), sessionId,
		"userId", authentication.GetInfo().GetUserId(),
		"userName", authentication.GetInfo().GetUserName()).Err()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = svc.RedisClient.Expire(ctx, sessionId, time.Minute*5).Err()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set sessionID cookie
	ctx.SetCookie("session_id", sessionId, 300, "", "", false, false)

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Log in succeeded"})
}

func (svc *WebService) CreateUser(ctx *gin.Context) {
	// Validate request
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
	_, _, userName, err := svc.CheckSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate request
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

// GetUserFollower gets followers of any user
func (svc *WebService) GetUserFollower(ctx *gin.Context) {
	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call GetUserFollower gprc service
	userFollower, err := svc.AuthenticateAndPostClient.GetUserFollower(ctx, &pb_aap.UserInfo{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Get user follower failed: %v", err)})
		return
	}

	// Return necessary information
	var followers []map[string]interface{}
	for _, follower := range userFollower.GetFollowers() {
		followers = append(followers, map[string]interface{}{"id": follower.UserId, "username": follower.UserName})
	}

	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "Get followers succeeded", "followers": followers})
}

func (svc *WebService) FollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.CheckSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call FollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.FollowUser(ctx,
		&pb_aap.UserAndFollower{
			User:     &pb_aap.UserInfo{UserId: int64(userId)},
			Follower: &pb_aap.UserInfo{UserId: int64(followerId)},
		})
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) UnfollowUser(ctx *gin.Context) {
	// Check sessionId authentication
	_, followerId, _, err := svc.CheckSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// Validate parameter
	stringUserId := ctx.Param("user_id")
	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID is not existed"})
		return
	}

	// Call UnfollowUser grpc service
	_, err = svc.AuthenticateAndPostClient.UnfollowUser(ctx,
		&pb_aap.UserAndFollower{
			User:     &pb_aap.UserInfo{UserId: int64(userId)},
			Follower: &pb_aap.UserInfo{UserId: int64(followerId)},
		},
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (svc *WebService) CheckSessionAuthentication(ctx *gin.Context) (sessionId string, userId int, userName string, err error) {
	sessionId, err = ctx.Cookie("session_id")
	if err != nil {
		return "", 0, "", err
	}

	userId, err = svc.RedisClient.HGet(svc.RedisClient.Context(), sessionId, "userId").Int()
	if err != nil {
		return "", 0, "", err
	}

	userName, err = svc.RedisClient.HGet(svc.RedisClient.Context(), sessionId, "userName").Result()
	if err != nil {
		return "", 0, "", err
	}

	return sessionId, userId, userName, nil
}

func (svc *WebService) ShouldBindAndValidateJSON(ctx *gin.Context, jsonRequest interface{}) error {
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		return err
	}

	err = validate.Struct(jsonRequest)
	if err != nil {
		return err
	}

	return nil
}
