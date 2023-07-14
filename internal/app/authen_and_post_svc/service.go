package authen_and_post_svc

import (
	"errors"
	"log"
	"os"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db          *gorm.DB
	kafkaWriter *kafka.Writer
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

	// Connect to kafka
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
	})
	if kafkaWriter == nil {
		return nil, errors.New("failed connecting to kafka writer")
	}

	return &AuthenticateAndPostService{db: db, kafkaWriter: kafkaWriter}, nil
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

// checkUserId checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) checkUserId(userId int64) (bool, types.User) {
	var userModel = types.User{}
	a.db.Raw("select * from user where id = ?", userId).Scan(&userModel)

	if userModel.ID == 0 {
		return false, types.User{}
	}
	return true, userModel
}

// checkPostId checks if an user with provided userId exists in database
func (a *AuthenticateAndPostService) checkPostId(postId int64) (bool, types.Post) {
	var postModel = types.Post{}
	a.db.Raw("select * from `post` where id = ?", postId).Scan(&postModel)

	if postModel.ID == 0 {
		return false, types.Post{}
	}
	return true, postModel
}

func (a *AuthenticateAndPostService) NewUserResult(userModel types.User) *pb_aap.UserResult {
	return &pb_aap.UserResult{
		Status: pb_aap.UserStatus_OK,
		Info: &pb_aap.UserDetailInfo{
			Id:           userModel.ID,
			UserName:     userModel.UserName,
			UserPassword: "",
			FirstName:    userModel.FirstName,
			LastName:     userModel.LastName,
			Dob:          userModel.DOB.Unix(),
			Email:        userModel.Email,
		},
	}
}

func (a *AuthenticateAndPostService) NewActionResult(status pb_aap.ActionStatus) *pb_aap.ActionResult {
	return &pb_aap.ActionResult{Status: status}
}
