package authen_and_post_svc

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	"github.com/maxuanquang/social-network/internal/utils"
	client_nfp "github.com/maxuanquang/social-network/pkg/client/newsfeed_publishing"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nfp "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AuthenticateAndPostService struct {
	pb_aap.UnimplementedAuthenticateAndPostServer
	db          *gorm.DB
	nfPubClient pb_nfp.NewsfeedPublishingClient
	redisClient *redis.Client

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

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
		})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &AuthenticateAndPostService{
		db:          db,
		nfPubClient: nfPubClient,
		redisClient: redisClient,
		logger:      logger,
	}, nil
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
