package authen_and_post_svc

import (
	"context"
	"errors"

	"github.com/maxuanquang/social-network/internal/auth"
	"github.com/maxuanquang/social-network/internal/models"
	aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"gorm.io/gorm"
)

func (a *AuthenticateAndPostService) CheckUserAuthentication(ctx context.Context, info *aap.UserInfo) (*aap.UserResult, error) {
	// Find user in the database using the provided information
	userModel := models.User{}

	err := a.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &aap.UserResult{}, err
	}

	err = auth.CheckPasswordHash(userModel.HashedPassword, info.GetUserPassword())
	if err != nil {
		return &aap.UserResult{}, err
	}

	// If the username and password are correct, return the necessary user information
	returnUserResult := aap.UserResult{
		Status: 1,
		Info: &aap.UserDetailInfo{
			UserId:    userModel.ID,
			UserName:  userModel.UserName,
			FirstName: userModel.FirstName,
			LastName:  userModel.LastName,
			Dob:       userModel.DOB.Unix(),
			Email:     userModel.Email,
		}}

	return &returnUserResult, err
}

func (a *AuthenticateAndPostService) CreateUser(ctx context.Context, info *aap.UserDetailInfo) (*aap.UserResult, error) {
	// Find user in the database using the provided information
	err := a.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&models.User{}).Error

	// If the user already exists, return an error
	if err == nil {
		return &aap.UserResult{}, errors.New("user already exists")
	}

	// If the user does not exist, create a new user
	salt, err := auth.GenerateRandomSalt()
	if err != nil {
		return &aap.UserResult{}, err
	}

	hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
	if err != nil {
		return &aap.UserResult{}, err
	}

	err = a.db.Exec("insert into user (id, hashed_password, salt, first_name, last_name, dob, email, user_name) values (null, ?, ?, ?, ?, ?, ?, ?)", hashed_password, string(salt), info.GetFirstName(), info.GetLastName(), info.GetDob(), info.GetEmail(), info.GetUserName()).Error
	if err != nil {
		return &aap.UserResult{}, err
	}

	returnUserResult := aap.UserResult{Status: 1, Info: info}
	return &returnUserResult, nil
}

func (a *AuthenticateAndPostService) EditUser(ctx context.Context, info *aap.UserDetailInfo) (*aap.UserResult, error) {
	// Check if the provided user exists in the database
	userModel := models.User{}
	err := a.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &aap.UserResult{}, err
	}

	// If the user exists, edit the information and return
	// Edit password
	if info.GetUserPassword() != "" {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return &aap.UserResult{}, err
		}

		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return &aap.UserResult{}, err
		}

		err = a.db.Exec("update user set hashed_password = ?, salt = ? where id = ?", hashed_password, salt, info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	// Edit first_name
	if info.GetFirstName() != "" {
		err = a.db.Exec("update user set first_name = ? where id = ?", info.GetFirstName(), info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	// Edit last_name
	if info.GetLastName() != "" {
		err = a.db.Exec("update user set last_name = ? where id = ?", info.GetLastName(), info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	// Edit dob
	if info.GetDob() != 0 {
		err = a.db.Exec("update user set dob = FROM_UNIXTIME(?) where id = ?", info.GetDob(), info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	// Edit email
	if info.GetEmail() != "" {
		err = a.db.Exec("update user set email = ? where id = ?", info.GetEmail(), info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	// Edit user_name
	if info.GetUserName() != "" {
		err = a.db.Exec("update user set username = ? where id = ?", info.GetUserName(), info.GetUserId()).Error
		if err != nil {
			return &aap.UserResult{}, err
		}
	}

	return &aap.UserResult{Status: 1, Info: info}, nil
}

func (a *AuthenticateAndPostService) GetUserFollower(ctx context.Context, info *aap.UserInfo) (*aap.UserFollower, error) {
	// Check if the user exists
	userModel := models.User{}
	err := a.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &aap.UserFollower{}, err
	}

	// If the user exists, return the followers
	var followers []models.User
	err = a.db.Raw("select follower_id from user_user where user_id = ?", userModel.ID).Scan(&followers).Error
	if err != nil {
		return &aap.UserFollower{}, err
	}

	returnUserFolower := aap.UserFollower{}
	for _, follower := range followers {
		followerInfo := aap.UserInfo{UserId: follower.ID, UserName: follower.UserName}
		returnUserFolower.Followers = append(returnUserFolower.Followers, &followerInfo)
	}

	return &returnUserFolower, nil
}

func (a *AuthenticateAndPostService) GetPostDetail(ctx context.Context, request *aap.GetPostRequest) (*aap.Post, error) {
	// Check if the post exists
	postModel := models.Post{}
	err := a.db.Raw("select * from post where id = ?", request.GetPostId()).Scan(&postModel).Error
	if err != nil {
		return &aap.Post{}, err
	}

	// If the post exists, return the post
	returnPost := aap.Post{
		PostId: postModel.ID,
		UserId: postModel.UserID, 
		ContentText: postModel.ContentText, 
		ContentImagePath: postModel.ContentImagePath, 
		Visible: postModel.Visible, 
		CreatedAt: postModel.CreatedAt.Unix(),
	}
	return &returnPost, nil
}

type AuthenticateAndPostService struct {
	aap.UnimplementedAuthenticateAndPostServer
	db *gorm.DB
}

func NewAuthenticateAndPostService(db *gorm.DB) *AuthenticateAndPostService {
	return &AuthenticateAndPostService{db: db}
}
