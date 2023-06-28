package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxuanquang/social-network/internal/pkg/types"

	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

// CheckUserNamePassword godoc
//
//	@Summary		get user
//	@Description	Check user's username and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request body types.LoginRequest true "login param"
//	@Success		200	{object} types.MessageResponse
//	@Failure		400	{object} types.MessageResponse
//	@Failure		500	{object} types.MessageResponse
//	@Router			/users/login [post]
func (svc *WebService) CheckUserAuthentication(ctx *gin.Context) {
	// Validate request
	var jsonRequest types.LoginRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}

	// Call CheckUserAuthentication service
	authentication, err := svc.AuthenticateAndPostClient.CheckUserAuthentication(ctx, &pb_aap.UserInfo{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// If logged in, set a sessionId for this session
	sessionId := uuid.New().String()

	// Save current sessionID and expiration time in Redis
	err = svc.RedisClient.HSet(svc.RedisClient.Context(), sessionId,
		"userId", authentication.GetInfo().GetUserId(),
		"userName", authentication.GetInfo().GetUserName()).Err()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	err = svc.RedisClient.Expire(ctx, sessionId, time.Minute*5).Err()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	// Set sessionID cookie
	ctx.SetCookie("session_id", sessionId, 300, "", "", false, false)

	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) CreateUser(ctx *gin.Context) {
	// Validate request
	var jsonRequest types.CreateUserRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
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
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "Failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
}

func (svc *WebService) EditUser(ctx *gin.Context) {
	// Check authorization
	_, _, userName, err := svc.checkSessionAuthentication(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, types.MessageResponse{Message: err.Error()})
		return
	}

	// Validate request
	var jsonRequest types.EditUserRequest
	err = ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
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
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusAccepted, types.MessageResponse{Message: "OK"})
}
