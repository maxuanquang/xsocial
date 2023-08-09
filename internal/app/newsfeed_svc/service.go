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
	newsfeedKey := "newsfeed:" + fmt.Sprint(request.GetUserId())
	postsIds, err := svc.redisClient.LPopCount(svc.redisClient.Context(), newsfeedKey, 5).Result()
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

// func (svc *NewsfeedService) getPostFromRedis(ctx context.Context, key string) (*types.RedisPost, error) {
// 	mapRedisPost, err := svc.redisClient.HGetAll(ctx, key).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var redisPost types.RedisPost
// 	err = svc.unmarshal(mapRedisPost, &redisPost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &redisPost, nil
// }

// func (svc *NewsfeedService) getCommentFromRedis(ctx context.Context, key string) (*types.RedisComment, error) {
// 	mapRedisComment, err := svc.redisClient.HGetAll(ctx, key).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var redisComment types.RedisComment
// 	err = svc.unmarshal(mapRedisComment, &redisComment)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &redisComment, nil
// }

// // unmarshal converts map[string]string to a struct.
// // It takes name of each field in map[string]string and maps with json tags in target struct
// func (svc *NewsfeedService) unmarshal(sourceMap map[string]string, objectPointer interface{}) error {
// 	objValue := reflect.ValueOf(objectPointer)
// 	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
// 		return fmt.Errorf("obj must be a non-nil pointer to a struct")
// 	}

// 	// Iterate over struct fields
// 	for i := 0; i < objValue.Elem().NumField(); i++ {
// 		field := objValue.Elem().Field(i)

// 		jsonTag := objValue.Elem().Type().Field(i).Tag.Get("json")
// 		mapValue, ok := sourceMap[jsonTag]
// 		if !ok || len(mapValue) == 0 {
// 			continue
// 		}

// 		switch field.Kind() {
// 		case reflect.Int64:
// 			intValue, err := strconv.ParseInt(mapValue, 10, 64)
// 			if err != nil {
// 				return err
// 			}
// 			field.SetInt(intValue)
// 		case reflect.String:
// 			field.SetString(mapValue)
// 		case reflect.Bool:
// 			boolValue, err := strconv.ParseBool(mapValue)
// 			if err != nil {
// 				return err
// 			}
// 			field.SetBool(boolValue)
// 		}
// 	}
// 	return nil
// }
