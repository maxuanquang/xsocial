package newsfeed_svc

import (
	"context"
	// "encoding/json"
	"errors"
	"fmt"

	// "reflect"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/utils"
	"go.uber.org/zap"

	// "github.com/maxuanquang/social-network/internal/pkg/types"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

type NewsfeedService struct {
	pb_nf.UnimplementedNewsfeedServer
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewNewsfeedService(cfg *configs.NewsfeedConfig) (*NewsfeedService, error) {
	// Connect to redisClient
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	return &NewsfeedService{
		redisClient: redisClient,
		logger:      logger,
	}, nil
}

func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, request *pb_nf.GetNewsfeedRequest) (*pb_nf.GetNewsfeedResponse, error) {
	// Query newsfeed from redis
	newsfeedKey := fmt.Sprintf("newsfeed:%d", request.GetUserId())
	postsIds, err := svc.redisClient.ZRevRange(svc.redisClient.Context(), newsfeedKey, 0, -1).Result()
	if errors.Is(err, redis.Nil) {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var int64PostsIds []int64
	for _, id := range postsIds {
		intPostId, err := strconv.Atoi(id)
		if err != nil {
			svc.logger.Debug(err.Error())
			continue
		}
		int64PostsIds = append(int64PostsIds, int64(intPostId))
	}
	if len(int64PostsIds) == 0 {
		return &pb_nf.GetNewsfeedResponse{
			Status: pb_nf.GetNewsfeedResponse_NEWSFEED_EMPTY,
		}, nil
	}
	return &pb_nf.GetNewsfeedResponse{
		Status:   pb_nf.GetNewsfeedResponse_OK,
		PostsIds: int64PostsIds,
	}, nil
}
