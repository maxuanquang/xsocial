package authen_and_post_svc

import (
	"context"
	"errors"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/auth"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	// Check if the userId which is changing information exists in the database
	existed, userModel := a.checkUserName(info.GetUserName())
	if !existed {
		return &pb_aap.UserResult{}, errors.New("user does not exist")
	}

	// If the user exists, edit the information and return
	var err error

	// Edit password
	if info.GetUserPassword() != "" {
		salt, err := auth.GenerateRandomSalt()
		if err != nil {
			return &pb_aap.UserResult{}, err
		}

		hashed_password, err := auth.HashPassword(info.GetUserPassword(), salt)
		if err != nil {
			return &pb_aap.UserResult{}, err
		}

		err = a.db.Exec("update user set hashed_password = ?, salt = ? where id = ?", hashed_password, salt, userModel.ID).Error
		if err != nil {
			return &pb_aap.UserResult{}, err
		}
	}

	// Edit first_name
	if info.GetFirstName() != "" {
		err = a.db.Exec("update user set first_name = ? where id = ?", info.GetFirstName(), userModel.ID).Error
		if err != nil {
			return &pb_aap.UserResult{}, err
		}
	}

	// Edit last_name
	if info.GetLastName() != "" {
		err = a.db.Exec("update user set last_name = ? where id = ?", info.GetLastName(), userModel.ID).Error
		if err != nil {
			return &pb_aap.UserResult{}, err
		}
	}

	// Edit dob
	if info.GetDob() >= -2208988800 { // 1900-01-01
		err = a.db.Exec("update user set dob = FROM_UNIXTIME(?) where id = ?", info.GetDob(), userModel.ID).Error
		if err != nil {
			return &pb_aap.UserResult{}, err
		}
	}

	// Edit email
	if info.GetEmail() != "" {
		err = a.db.Exec("update user set email = ? where id = ?", info.GetEmail(), userModel.ID).Error
		if err != nil {
			return &pb_aap.UserResult{}, err
		}
	}

	// Return the necessary user information
	_, userModel = a.checkUserName(info.GetUserName())
	return a.NewUserResult(userModel), nil
}

func (a *AuthenticateAndPostService) GetUserFollower(ctx context.Context, info *pb_aap.UserInfo) (*pb_aap.UserFollower, error) {
	// Check if the user exists
	userModel := types.User{}
	err := a.db.Raw("select * from user where user_name = ?", info.GetUserName()).Scan(&userModel).Error
	if err != nil {
		return &pb_aap.UserFollower{}, err
	}

	// If the user exists, return the followers
	var followers []types.User
	err = a.db.Raw("select follower_id from user_user where user_id = ?", userModel.ID).Scan(&followers).Error
	if err != nil {
		return &pb_aap.UserFollower{}, err
	}

	returnUserFolower := pb_aap.UserFollower{}
	for _, follower := range followers {
		followerInfo := pb_aap.UserInfo{UserId: follower.ID, UserName: follower.UserName}
		returnUserFolower.Followers = append(returnUserFolower.Followers, &followerInfo)
	}

	return &returnUserFolower, nil
}

func (a *AuthenticateAndPostService) GetPostDetail(ctx context.Context, request *pb_aap.GetPostRequest) (*pb_aap.Post, error) {
	// Check if the post exists
	postModel := types.Post{}
	err := a.db.Raw("select * from post where id = ?", request.GetPostId()).Scan(&postModel).Error
	if err != nil {
		return &pb_aap.Post{}, err
	}

	// If the post exists, return the post
	returnPost := pb_aap.Post{
		PostId:           postModel.ID,
		UserId:           postModel.UserID,
		ContentText:      postModel.ContentText,
		ContentImagePath: postModel.ContentImagePath,
		Visible:          postModel.Visible,
		CreatedAt:        postModel.CreatedAt.Unix(),
	}
	return &returnPost, nil
}

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db *gorm.DB
}

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Connect to database
	mysqlConfig := mysql.Config{
		DSN: cfg.MySQL.DSN,
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		return &AuthenticateAndPostService{}, err
	}

	return &AuthenticateAndPostService{db: db}, err
}

// checkUserName checks if an user with provided username exists in database
func (a *AuthenticateAndPostService) checkUserName(username string) (bool, types.User) {
	var userModel = types.User{}
	a.db.Raw("select * from user where user_name = ?", username).Scan(&userModel)

	if userModel.ID == 0 {
		return false, types.User{}
	}
	return true, userModel
}

func (a *AuthenticateAndPostService) NewUserResult(userModel types.User) *pb_aap.UserResult {
	return &pb_aap.UserResult{
		Status: pb_aap.UserStatus_OK,
		Info: &pb_aap.UserDetailInfo{
			UserId:       userModel.ID,
			UserName:     userModel.UserName,
			UserPassword: "",
			FirstName:    userModel.FirstName,
			LastName:     userModel.LastName,
			Dob:          userModel.DOB.Unix(),
			Email:        userModel.Email,
		},
	}
}
