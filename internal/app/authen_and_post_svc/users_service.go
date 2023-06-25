package authen_and_post_svc

import (
	"context"
	"errors"
	"github.com/maxuanquang/social-network/internal/auth"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
)

func (a *AuthenticateAndPostService) CreateUser(ctx context.Context, info *pb_aap.UserDetailInfo) (*pb_aap.UserResult, error) {
	existed, _ := a.checkUserName(info.GetUserName())
	if existed {
		return nil, errors.New("user already exist")
	}

	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		return nil, err
	}

	hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return nil, err
	}

	err = a.db.Exec(
		"insert into user (id, hashed_password, salt, first_name, last_name, dob, email, user_name) values (null, ?, ?, ?, ?, FROM_UNIXTIME(?), ?, ?)",
		hashed_password,
		salt,
		info.GetFirstName(),
		info.GetLastName(),
		info.GetDob(),
		info.GetEmail(),
		info.GetUserName(),
	).Error
	if err != nil {
		return nil, err
	}

	// Return the necessary user information
	_, userModel := a.checkUserName(info.GetUserName())
	return a.NewUserResult(userModel), nil
}

func (a *AuthenticateAndPostService) CheckUserAuthentication(ctx context.Context, info *pb_aap.UserInfo) (*pb_aap.UserResult, error) {
	// Find user in the database using the provided information
	existed, userModel := a.checkUserName(info.GetUserName())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// Check password matching
	err := auth.CheckPasswordHash(userModel.HashedPassword, info.GetUserPassword(), userModel.Salt)
	if err != nil {
		return nil, err
	}

	// Return the necessary user information
	return a.NewUserResult(userModel), nil
}

func (a *AuthenticateAndPostService) EditUser(ctx context.Context, info *pb_aap.UserDetailInfo) (*pb_aap.UserResult, error) {
	// Check if the username which is changing information exists in the database
	existed, userModel := a.checkUserName(info.GetUserName())
	if !existed {
		return nil, errors.New("user does not exist")
	}

	// If the user exists, edit the information and return
	var err error

	// Edit password
	if info.GetUserPassword() != "" {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return nil, err
		}

		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return nil, err
		}

		err = a.db.Exec("update user set hashed_password = ?, salt = ? where id = ?", hashed_password, salt, userModel.ID).Error
		if err != nil {
			return nil, err
		}
	}

	// Edit first_name
	if info.GetFirstName() != "" {
		err = a.db.Exec("update user set first_name = ? where id = ?", info.GetFirstName(), userModel.ID).Error
		if err != nil {
			return nil, err
		}
	}

	// Edit last_name
	if info.GetLastName() != "" {
		err = a.db.Exec("update user set last_name = ? where id = ?", info.GetLastName(), userModel.ID).Error
		if err != nil {
			return nil, err
		}
	}

	// Edit dob
	if info.GetDob() >= -2208988800 { // 1900-01-01
		err = a.db.Exec("update user set dob = FROM_UNIXTIME(?) where id = ?", info.GetDob(), userModel.ID).Error
		if err != nil {
			return nil, err
		}
	}

	// Edit email
	if info.GetEmail() != "" {
		err = a.db.Exec("update user set email = ? where id = ?", info.GetEmail(), userModel.ID).Error
		if err != nil {
			return nil, err
		}
	}

	// Return the necessary user information
	_, userModel = a.checkUserName(info.GetUserName())
	return a.NewUserResult(userModel), nil
}
