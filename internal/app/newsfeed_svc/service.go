package newsfeed_svc

import (
	"context"
	// "encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

type NewsfeedService struct {
	pb_nf.UnimplementedNewsfeedServer
	redisClient *redis.Client
}

func NewNewsfeedService(cfg *configs.NewsfeedConfig) (*NewsfeedService, error) {
	// Connect to redisClient
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	return &NewsfeedService{
		redisClient: redisClient,
	}, nil
}

func (svc *NewsfeedService) GetNewsfeed(ctx context.Context, request *pb_nf.NewsfeedRequest) (*pb_nf.NewsfeedResponse, error) {
	// Query newsfeed from redis
	newsfeedKey := "newsfeed-" + fmt.Sprint(request.UserId)
	postsIds, err := svc.redisClient.LPopCount(svc.redisClient.Context(), newsfeedKey, 5).Result()
	if err != nil {
		return nil, err
	}
	if len(postsIds) == 0 {
		return nil, errors.New("no new posts in newsfeed")
	}

	var posts []*pb_aap.PostDetailInfo
	for _, id := range postsIds {
		// Get post from redis
		postKey := "post-" + id
		redisPost, err := svc.getPostFromRedis(ctx, postKey)
		if err != nil {
			return nil, err
		}

		// Get post's comments from redis
		var comments []*pb_aap.CommentInfo
		if len(redisPost.CommentsIds) > 0 {
			for _, comment_id := range strings.Split(redisPost.CommentsIds, " ") {
				commentKey := "comment-" + comment_id
				redisComment, err := svc.getCommentFromRedis(ctx, commentKey)
				if err != nil {
					return nil, err
				}

				comments = append(comments, &pb_aap.CommentInfo{
					Id:      redisComment.ID,
					PostId:  redisComment.PostID,
					UserId:  redisComment.UserID,
					Content: redisComment.Content,
				})
			}
		}

		// Get post's likes from redis
		var likes []*pb_aap.LikeInfo
		if len(redisPost.LikedUsersIds) > 0 {
			for _, liked_user_id := range strings.Split(redisPost.LikedUsersIds, " ") {
				user_id, err := strconv.Atoi(liked_user_id)
				if err != nil {
					return nil, err
				}
				likes = append(likes, &pb_aap.LikeInfo{
					UserId: int64(user_id),
					PostId: redisPost.ID,
				})
			}
		}

		posts = append(posts, &pb_aap.PostDetailInfo{
			Id:               redisPost.ID,
			UserId:           redisPost.UserID,
			ContentText:      redisPost.ContentText,
			ContentImagePath: strings.Split(redisPost.ContentImagePath, " "),
			Visible:          redisPost.Visible,
			CreatedAt:        redisPost.CreatedAt,
			Comments:         comments,
			Likes:            likes,
		})

	}
	return &pb_nf.NewsfeedResponse{Posts: posts}, nil
}

func (svc *NewsfeedService) getPostFromRedis(ctx context.Context, key string) (*types.RedisPost, error) {
	mapRedisPost, err := svc.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var redisPost types.RedisPost
	err = svc.unmarshal(mapRedisPost, &redisPost)
	if err != nil {
		return nil, err
	}

	return &redisPost, nil
}

func (svc *NewsfeedService) getCommentFromRedis(ctx context.Context, key string) (*types.RedisComment, error) {
	mapRedisComment, err := svc.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var redisComment types.RedisComment
	err = svc.unmarshal(mapRedisComment, &redisComment)
	if err != nil {
		return nil, err
	}
	return &redisComment, nil
}

// unmarshal converts map[string]string to a struct.
// It takes name of each field in map[string]string and maps with json tags in target struct
func (svc *NewsfeedService) unmarshal(sourceMap map[string]string, objectPointer interface{}) error {
	objValue := reflect.ValueOf(objectPointer)
	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
		return fmt.Errorf("obj must be a non-nil pointer to a struct")
	}

	// Iterate over struct fields
	for i := 0; i < objValue.Elem().NumField(); i++ {
		field := objValue.Elem().Field(i)

		jsonTag := objValue.Elem().Type().Field(i).Tag.Get("json")
		mapValue, ok := sourceMap[jsonTag]
		if !ok || len(mapValue) == 0 {
			continue
		}

		switch field.Kind() {
		case reflect.Int64:
			intValue, err := strconv.ParseInt(mapValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.String:
			field.SetString(mapValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(mapValue)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		}
	}
	return nil
}

// func (svc *NewsfeedService) convertMapToPost(mapRedisPost map[string]string) (*types.RedisPost, error) {
// 	redisPost := &types.RedisPost{}
// 	redisPostType := reflect.TypeOf(redisPost).Elem()
// 	redisPostValue := reflect.ValueOf(redisPost).Elem()

// 	// Iterate over struct fields
// 	for i := 0; i < redisPostType.NumField(); i++ {
// 		jsonTag := redisPostType.Field(i).Tag.Get("json")

// 		value, ok := mapRedisPost[jsonTag]
// 		if !ok {
// 			continue
// 		}

// 		if redisPostType.Field(i).Type.Kind() == reflect.Int64 {
// 			valueTemp, err := strconv.Atoi(value)
// 			if err != nil {
// 				return nil, err
// 			}
// 			valueConverted := reflect.ValueOf(valueTemp).Convert(redisPostType.Field(i).Type)
// 			redisPostValue.Field(i).Set(valueConverted)
// 		} else if redisPostType.Field(i).Type.Kind() == reflect.String {
// 			valueTemp := value
// 			valueConverted := reflect.ValueOf(valueTemp).Convert(redisPostType.Field(i).Type)
// 			redisPostValue.Field(i).Set(valueConverted)
// 		} else if redisPostType.Field(i).Type.Kind() == reflect.Bool {
// 			var valueTemp bool
// 			if value == "1" || value == "true" {
// 				valueTemp = true
// 			} else {
// 				valueTemp = false
// 			}
// 			valueConverted := reflect.ValueOf(valueTemp).Convert(redisPostType.Field(i).Type)
// 			redisPostValue.Field(i).Set(valueConverted)
// 		}
// 	}
// 	return redisPost, nil
// }
