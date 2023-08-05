package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	authentication, err := svc.AuthenticateAndPostClient.CheckUserAuthentication(ctx, &pb_aap.CheckUserAuthenticationRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "wrong username or password"})
		return
	} else if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_WRONG_PASSWORD {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "wrong username or password"})
		return
	} else if authentication.GetStatus() == pb_aap.CheckUserAuthenticationResponse_OK {
		// Set a sessionId for this session
		sessionId := uuid.New().String()

		// Save current sessionID and expiration time in Redis
		err = svc.RedisClient.Set(svc.RedisClient.Context(), sessionId, authentication.GetUserId(), time.Minute*15).Err()
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
			return
		}

		// Set sessionID cookie
		ctx.SetCookie("session_id", sessionId, 900, "", "", false, false)

		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
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
	dob, err := time.Parse(time.DateOnly, jsonRequest.DateOfBirth)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
		return
	}
	resp, err := svc.AuthenticateAndPostClient.CreateUser(ctx, &pb_aap.CreateUserRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
		FirstName:    jsonRequest.FirstName,
		LastName:     jsonRequest.LastName,
		DateOfBirth:  timestamppb.New(dob),
		Email:        jsonRequest.Email,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CreateUserResponse_USERNAME_EXISTED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "username existed"})
		return
	} else if resp.GetStatus() == pb_aap.CreateUserResponse_EMAIL_EXISTED {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "email existed"})
		return
	} else if resp.GetStatus() == pb_aap.CreateUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) EditUser(ctx *gin.Context) {
	// Check authorization
	_, userId, err := svc.checkSessionAuthentication(ctx)
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
	var password *string
	if jsonRequest.Password != "" {
		password = &jsonRequest.Password
	}
	var firstName *string
	if jsonRequest.FirstName != "" {
		firstName = &jsonRequest.FirstName
	}
	var lastName *string
	if jsonRequest.LastName != "" {
		lastName = &jsonRequest.LastName
	}
	var dateOfBirth *timestamppb.Timestamp
	if jsonRequest.DateOfBirth != "" {
		dob, err := time.Parse(time.DateOnly, jsonRequest.DateOfBirth)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
			return
		}
		dateOfBirth = timestamppb.New(dob)
	}

	// Call EditUser service
	resp, err := svc.AuthenticateAndPostClient.EditUser(ctx, &pb_aap.EditUserRequest{
		UserId:       int64(userId),
		UserPassword: password,
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.EditUserResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.EditUserResponse_OK {
		ctx.IndentedJSON(http.StatusOK, types.MessageResponse{Message: "OK"})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) GetUserDetailInfo(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call gprc service
	resp, err := svc.AuthenticateAndPostClient.GetUserDetailInfo(ctx, &pb_aap.GetUserDetailInfoRequest{
		UserId: int64(userId),
	})
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.GetUserDetailInfoResponse_USER_NOT_FOUND {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "user not found"})
		return
	} else if resp.GetStatus() == pb_aap.GetUserDetailInfoResponse_OK {
		ctx.IndentedJSON(http.StatusAccepted, gin.H{
			"user_id":       resp.GetUser().GetUserId(),
			"user_name":     resp.GetUser().GetUserName(),
			"first_name":    resp.GetUser().GetFirstName(),
			"last_name":     resp.GetUser().GetLastName(),
			"date_of_birth": resp.GetUser().GetDateOfBirth().AsTime().Format(time.DateOnly),
			"email":         resp.GetUser().GetEmail(),
		})
		return
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, types.MessageResponse{Message: "unknown error"})
		return
	}
}

func (svc *WebService) checkSessionAuthentication(ctx *gin.Context) (sessionId string, userId int, err error) {
	sessionId, err = ctx.Cookie("session_id")
	if err != nil {
		return "", 0, err
	}

	userId, err = svc.RedisClient.Get(svc.RedisClient.Context(), sessionId).Int()
	if err != nil {
		return "", 0, err
	}

	return sessionId, userId, nil
}
