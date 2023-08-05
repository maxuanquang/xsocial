package authen_and_post_svc

import (
	"encoding/json"
	"errors"
	// "log"
	// "os"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	client_nfp "github.com/maxuanquang/social-network/pkg/client/newsfeed_publishing"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nfp "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	// "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db          *gorm.DB
	nfPubClient pb_nfp.NewsfeedPublishingClient

	logger *zap.Logger
}

func NewAuthenticateAndPostService(cfg *configs.AuthenticateAndPostConfig) (*AuthenticateAndPostService, error) {
	// Connect to database
	mysqlConfig := mysql.Config{
		DSN: cfg.MySQL.DSN,
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Connect to NewsfeedPublishingClient
	nfPubClient, err := client_nfp.NewClient(cfg.NewsfeedPublishing.Hosts)
	if err != nil {
		return nil, err
	}

	// Establish logger
	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	return &AuthenticateAndPostService{
		db:          db,
		nfPubClient: nfPubClient,
		logger:      logger,
	}, nil
}

func newLogger() (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {"foo": "bar"},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	logger := zap.Must(cfg.Build())
	return logger, nil
}

// findUserById checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) findUserById(userId int64) (exist bool, user types.User) {
	result := a.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

// findUserByUserName checks if an user with provided username exists in database
func (a *AuthenticateAndPostService) findUserByUserName(userName string) (exist bool, user types.User) {
	result := a.db.Where(&types.User{UserName: userName}).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.User{}
	}
	return true, user
}

// // findUserByEmail checks if an user with provided username exists in database
// func (a *AuthenticateAndPostService) findUserByEmail(email string) (exist bool, user types.User) {
// 	result := a.db.Where(&types.User{Email: email}).First(&user)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return false, types.User{}
// 	}
// 	return true, user
// }

// findPostById checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) findPostById(postId int64) (exist bool, post types.Post) {
	result := a.db.First(&post, postId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, types.Post{}
	}
	return true, post
}

// func (a *AuthenticateAndPostService) NewUserResult(userModel types.User) *pb_aap.UserResult {
// 	return &pb_aap.UserResult{
// 		Status: pb_aap.UserStatus_OK,
// 		Info: &pb_aap.UserDetailInfo{
// 			Id:           userModel.ID,
// 			UserName:     userModel.UserName,
// 			UserPassword: "",
// 			FirstName:    userModel.FirstName,
// 			LastName:     userModel.LastName,
// 			Dob:          userModel.DOB.Unix(),
// 			Email:        userModel.Email,
// 		},
// 	}
// }

// func (a *AuthenticateAndPostService) NewActionResult(status pb_aap.ActionStatus) *pb_aap.ActionResult {
// 	return &pb_aap.ActionResult{Status: status}
// }
