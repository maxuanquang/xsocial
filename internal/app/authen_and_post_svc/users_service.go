package authen_and_post_svc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maxuanquang/social-network/internal/auth"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (a *AuthenticateAndPostService) CreateUser(ctx context.Context, info *pb_aap.CreateUserRequest) (*pb_aap.CreateUserResponse, error) {
	// Check user name existence
	exist, _ := a.findUserByUserName(info.GetUserName())
	if exist {
		return &pb_aap.CreateUserResponse{Status: pb_aap.CreateUserResponse_USERNAME_EXISTED}, nil
	}

	// Password hash and salt
	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		return nil, err
	}
	hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return nil, err
	}

	// Create user
	newUser := types.User{
		HashedPassword: hashed_password,
		Salt:           salt,
		Email:          info.GetEmail(),
		UserName:       info.GetUserName(),
	}
	result := a.db.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb_aap.CreateUserResponse{
		Status: pb_aap.CreateUserResponse_OK,
		UserId: int64(newUser.ID),
	}, nil
}

func (a *AuthenticateAndPostService) CheckUserAuthentication(ctx context.Context, info *pb_aap.CheckUserAuthenticationRequest) (*pb_aap.CheckUserAuthenticationResponse, error) {
	// Check user name
	var user types.User
	result := a.db.Where(&types.User{UserName: info.GetUserName()}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &pb_aap.CheckUserAuthenticationResponse{Status: pb_aap.CheckUserAuthenticationResponse_USER_NOT_FOUND}, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	// Password matching
	err := auth.CheckPasswordHash(user.HashedPassword, info.GetUserPassword(), user.Salt)
	if err != nil {
		return &pb_aap.CheckUserAuthenticationResponse{Status: pb_aap.CheckUserAuthenticationResponse_WRONG_PASSWORD}, nil
	}

	return &pb_aap.CheckUserAuthenticationResponse{
		Status: pb_aap.CheckUserAuthenticationResponse_OK,
		User: &pb_aap.UserDetailInfo{
			UserId:         int64(user.ID),
			UserName:       user.UserName,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			DateOfBirth:    timestamppb.New(user.DateOfBirth.Time),
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			CoverPicture:   user.CoverPicture,
		},
	}, nil
}

func (a *AuthenticateAndPostService) EditUser(ctx context.Context, info *pb_aap.EditUserRequest) (*pb_aap.EditUserResponse, error) {
	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.EditUserResponse{Status: pb_aap.EditUserResponse_USER_NOT_FOUND}, nil
	}
	if info.FirstName != nil {
		user.FirstName = info.GetFirstName()
	}
	if info.LastName != nil {
		user.LastName = info.GetLastName()
	}
	if info.DateOfBirth != nil {
		user.DateOfBirth = sql.NullTime{Time: info.GetDateOfBirth().AsTime()}
	}
	if info.UserPassword != nil {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return nil, err
		}
		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return nil, err
		}
		user.Salt = salt
		user.HashedPassword = hashed_password
	}
	if info.ProfilePicture != nil {
		user.ProfilePicture = info.GetProfilePicture()
	}
	if info.CoverPicture != nil {
		user.CoverPicture = info.GetCoverPicture()
	}
	a.db.Save(&user)

	return &pb_aap.EditUserResponse{
		Status: pb_aap.EditUserResponse_OK,
	}, nil
}

func (a *AuthenticateAndPostService) GetUserDetailInfo(ctx context.Context, info *pb_aap.GetUserDetailInfoRequest) (*pb_aap.GetUserDetailInfoResponse, error) {
	userKey := fmt.Sprintf("user:%d", info.GetUserId())
	cacheExist := (a.redisClient.Exists(context.Background(), userKey).Val() == 1)
	if cacheExist {
		userJson := a.redisClient.Get(context.Background(), userKey).Val()
		var user types.User
		err := json.Unmarshal([]byte(userJson), &user)
		if err == nil {
			return &pb_aap.GetUserDetailInfoResponse{
				Status: pb_aap.GetUserDetailInfoResponse_OK,
				User: &pb_aap.UserDetailInfo{
					UserId:      int64(user.ID),
					UserName:    user.UserName,
					FirstName:   user.FirstName,
					LastName:    user.LastName,
					DateOfBirth: timestamppb.New(user.DateOfBirth.Time),
					Email:       user.Email,
				},
			}, nil
		}
	}

	exist, user := a.findUserById(info.GetUserId())
	if !exist {
		return &pb_aap.GetUserDetailInfoResponse{Status: pb_aap.GetUserDetailInfoResponse_USER_NOT_FOUND}, nil
	}

	return &pb_aap.GetUserDetailInfoResponse{
		Status: pb_aap.GetUserDetailInfoResponse_OK,
		User: &pb_aap.UserDetailInfo{
			UserId:      int64(user.ID),
			UserName:    user.UserName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DateOfBirth: timestamppb.New(user.DateOfBirth.Time),
			Email:       user.Email,
		},
	}, nil
}
