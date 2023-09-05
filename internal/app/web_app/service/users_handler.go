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

// CheckUserNamePassword checks user's authentication
//
//	@Summary		Check user's username and password
//	@Description	check user's username and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.LoginRequest	true	"Login parameters"
//	@Success		200		{object}	types.LoginResponse
//	@Failure		400		{object}	types.LoginResponse
//	@Failure		500		{object}	types.LoginResponse
//	@Router			/users/login [post]
func (svc *WebService) CheckUserAuthentication(ctx *gin.Context) {
	// Metrics for Prometheus
	svc.countReporter.WithLabelValues("check_user_login", "total").Inc()
	var start = time.Now()
	var end time.Time
	var httpStatus int
	defer func() {
		svc.countReporter.WithLabelValues("check_user_login", strconv.Itoa(httpStatus)).Inc()
		svc.latencyReporter.WithLabelValues("check_user_login", strconv.Itoa(httpStatus)).Observe(float64(end.UnixMilli() - start.UnixMilli()))
	}()

	// Validate request
	var jsonRequest types.LoginRequest
	err := ctx.ShouldBindJSON(&jsonRequest)
	if err != nil {
		httpStatus = http.StatusBadRequest
		end = time.Now()
		ctx.IndentedJSON(http.StatusBadRequest, types.LoginResponse{Message: err.Error()})
		return
	}
	err = validate.Struct(jsonRequest)
	if err != nil {
		httpStatus = http.StatusBadRequest
		end = time.Now()
		ctx.IndentedJSON(http.StatusBadRequest, types.LoginResponse{Message: err.Error()})
		return
	}

	// Call CheckUserAuthentication service
	resp, err := svc.authenticateAndPostClient.CheckUserAuthentication(ctx, &pb_aap.CheckUserAuthenticationRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
	})
	if err != nil {
		httpStatus = http.StatusInternalServerError
		end = time.Now()
		ctx.IndentedJSON(http.StatusInternalServerError, types.LoginResponse{Message: err.Error()})
		return
	}
	if resp.GetStatus() == pb_aap.CheckUserAuthenticationResponse_USER_NOT_FOUND {
		httpStatus = http.StatusOK
		end = time.Now()
		ctx.IndentedJSON(http.StatusOK, types.LoginResponse{Message: "wrong username or password"})
		return
	} else if resp.GetStatus() == pb_aap.CheckUserAuthenticationResponse_WRONG_PASSWORD {
		httpStatus = http.StatusOK
		end = time.Now()
		ctx.IndentedJSON(http.StatusOK, types.LoginResponse{Message: "wrong username or password"})
		return
	} else if resp.GetStatus() == pb_aap.CheckUserAuthenticationResponse_OK {
		// Set a sessionId for this session
		sessionId := uuid.New().String()

		// Save current sessionID and expiration time in Redis
		svc.redisClient.Set(svc.redisClient.Context(), sessionId, resp.GetUser().GetUserId(), time.Minute*15)

		// Set sessionID cookie
		// (this cookie is currently not working on Google Chrome)
		// TODO: Add HTTPS for cookie to work properly
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "session_id",
			Value:    sessionId,
			MaxAge:   900,
			Path:     "/",
			Domain:   "",
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
			HttpOnly: false,
		})

		httpStatus = http.StatusOK
		end = time.Now()
		ctx.IndentedJSON(http.StatusOK, types.LoginResponse{
			Message: "OK",
			User: types.UserDetailInfo{
				UserID:         resp.GetUser().GetUserId(),
				UserName:       resp.GetUser().GetUserName(),
				FirstName:      resp.GetUser().GetFirstName(),
				LastName:       resp.GetUser().GetLastName(),
				DateOfBirth:    resp.GetUser().GetDateOfBirth().AsTime().Format(time.DateOnly),
				Email:          resp.GetUser().GetEmail(),
				ProfilePicture: resp.GetUser().GetProfilePicture(),
				CoverPicture:   resp.GetUser().GetCoverPicture(),
			}})
		return
	} else {
		httpStatus = http.StatusInternalServerError
		end = time.Now()
		ctx.IndentedJSON(http.StatusInternalServerError, types.LoginResponse{Message: "unknown error"})
		return
	}
}

// CreateUser creates new user account
//
//	@Summary		create new user account
//	@Description	create new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.CreateUserRequest	true	"Create user parameters"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/users/signup [post]
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
	resp, err := svc.authenticateAndPostClient.CreateUser(ctx, &pb_aap.CreateUserRequest{
		UserName:     jsonRequest.UserName,
		UserPassword: jsonRequest.Password,
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

// EditUser edits user information
//
//	@Summary		edit user information
//	@Description	edit user information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.EditUserRequest	true	"Edit user information parameters"
//	@Success		200		{object}	types.MessageResponse
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/users/edit [put]
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
	if jsonRequest.Password != nil {
		password = jsonRequest.Password
	}
	var firstName *string
	if jsonRequest.FirstName != nil {
		firstName = jsonRequest.FirstName
	}
	var lastName *string
	if jsonRequest.LastName != nil {
		lastName = jsonRequest.LastName
	}
	var dateOfBirth *timestamppb.Timestamp
	if jsonRequest.DateOfBirth != nil {
		dob, err := time.Parse(time.DateOnly, *jsonRequest.DateOfBirth)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: err.Error()})
			return
		}
		dateOfBirth = timestamppb.New(dob)
	}
	var profilePicture *string
	if jsonRequest.ProfilePicture != nil {
		profilePicture = jsonRequest.ProfilePicture
	}
	var coverPicture *string
	if jsonRequest.CoverPicture != nil {
		coverPicture = jsonRequest.CoverPicture
	}

	// Call EditUser service
	resp, err := svc.authenticateAndPostClient.EditUser(ctx, &pb_aap.EditUserRequest{
		UserId:         int64(userId),
		UserPassword:   password,
		FirstName:      firstName,
		LastName:       lastName,
		DateOfBirth:    dateOfBirth,
		ProfilePicture: profilePicture,
		CoverPicture:   coverPicture,
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

// GetUserDetailInfo gets user information
//
//	@Summary		get user information
//	@Description	get user information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		int	true	"User ID"
//	@Success		200		{object}	types.UserDetailInfo
//	@Failure		400		{object}	types.MessageResponse
//	@Failure		500		{object}	types.MessageResponse
//	@Router			/users/{user_id} [get]
func (svc *WebService) GetUserDetailInfo(ctx *gin.Context) {
	// Check URL params
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, types.MessageResponse{Message: "post not found"})
		return
	}

	// Call gprc service
	resp, err := svc.authenticateAndPostClient.GetUserDetailInfo(ctx, &pb_aap.GetUserDetailInfoRequest{
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
		ctx.IndentedJSON(http.StatusAccepted, types.UserDetailInfo{
			UserID:         resp.GetUser().GetUserId(),
			UserName:       resp.GetUser().GetUserName(),
			FirstName:      resp.GetUser().GetFirstName(),
			LastName:       resp.GetUser().GetLastName(),
			DateOfBirth:    resp.GetUser().GetDateOfBirth().AsTime().Format(time.DateOnly),
			Email:          resp.GetUser().GetEmail(),
			ProfilePicture: resp.GetUser().GetProfilePicture(),
			CoverPicture:   resp.GetUser().GetCoverPicture(),
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

	userId, err = svc.redisClient.Get(svc.redisClient.Context(), sessionId).Int()
	if err != nil {
		return "", 0, err
	}

	return sessionId, userId, nil
}
