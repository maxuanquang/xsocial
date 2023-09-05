package newsfeed_publishing_svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/utils"
	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nfp "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type NewsfeedPublishingService struct {
	pb_nfp.UnimplementedNewsfeedPublishingServer
	kafkaWriter               *kafka.Writer
	kafkaReader               *kafka.Reader
	redisClient               *redis.Client
	authenticateAndPostClient pb_aap.AuthenticateAndPostClient

	logger *zap.Logger
}

func NewNewsfeedPublishingService(cfg *configs.NewsfeedPublishingConfig) (*NewsfeedPublishingService, error) {
	// Connect to kafka writer
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka writer: ", 0),
		Async:   true,
	})
	if kafkaWriter == nil {
		return nil, errors.New("failed creating kafka writer")
	}

	// Connect to kafka reader
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka reader: ", 0),
		GroupID: "0",
	})
	if kafkaReader == nil {
		return nil, errors.New("failed creating kafka reader")
	}

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Connect to aap service
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// Establish logger
	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Return
	return &NewsfeedPublishingService{
		kafkaWriter:               kafkaWriter,
		kafkaReader:               kafkaReader,
		redisClient:               redisClient,
		authenticateAndPostClient: aapClient,
		logger:                    logger,
	}, nil
}

func (svc *NewsfeedPublishingService) PublishPost(ctx context.Context, info *pb_nfp.PublishPostRequest) (*pb_nfp.PublishPostResponse, error) {
	value := map[string]int64{
		"user_id":    info.GetUserId(),
		"post_id":    info.GetPostId(),
		"created_at": info.GetCreatedAt().GetSeconds(),
	}
	jsonValue, _ := json.Marshal(value)
	err := svc.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("post"),
		Value: jsonValue,
	})
	if err != nil {
		return nil, err
	}

	return &pb_nfp.PublishPostResponse{
		Status: pb_nfp.PublishPostResponse_OK,
	}, nil
}

func (svc *NewsfeedPublishingService) Run() {
	for {
		message, err := svc.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}
		svc.processMessage(message)
	}
}

func (svc *NewsfeedPublishingService) processMessage(message kafka.Message) {
	msgType := string(message.Key)

	// Process message based on its key
	if msgType == "post" {
		svc.processPost(message.Value)
	}
}

func (svc *NewsfeedPublishingService) processPost(value []byte) {
	var message map[string]int64
	err := json.Unmarshal(value, &message)
	if err != nil {
		panic(err)
	}

	// Find followers of user that created post
	followersKey := fmt.Sprintf("followers:%d", message["user_id"])
	exist := (svc.redisClient.Exists(context.Background(), followersKey).Val() == 1)
	if !exist {
		resp, err := svc.authenticateAndPostClient.GetUserFollower(
			context.Background(),
			&pb_aap.GetUserFollowerRequest{
				UserId: message["user_id"],
			})
		if err != nil {
			panic(err)
		}

		var followersIds []interface{}
		for _, id := range resp.GetFollowersIds() {
			followersIds = append(followersIds, id)
		}
		svc.redisClient.RPush(context.Background(), followersKey, followersIds...)
		svc.redisClient.Expire(context.Background(), followersKey, 15*time.Minute)
	}
	followersIds := svc.redisClient.LRange(context.Background(), followersKey, 0, -1).Val()

	// Add this post_id into followers' newsfeed
	for _, id := range followersIds {
		newsfeedKey := fmt.Sprintf("newsfeed:%s", id)
		svc.redisClient.ZAdd(context.Background(), newsfeedKey, &redis.Z{
			Score:  float64(message["created_at"]),
			Member: message["post_id"],
		})
		svc.redisClient.Expire(context.Background(), newsfeedKey, 15*time.Minute)
	}
}

// func (svc *NewsfeedPublishingService) acquireDistributedLock(lockName string) {
// 	for {
// 		err := svc.redisClient.SetNX(context.Background(), lockName, "lock", 5*time.Minute).Err()
// 		if err == nil {
// 			break
// 		}
// 		time.Sleep(time.Second / 2)
// 		continue
// 	}
// }

// func (svc *NewsfeedPublishingService) releaseDistributedLock(lockName string) {
// 	svc.redisClient.Del(context.Background(), lockName)
// }

// func (svc *NewsfeedPublishingService) newMapRedisPost(postDetailInfo *pb_aap.PostDetailInfo) map[string]interface{} {
// 	var commentsIds []string
// 	for _, comment := range postDetailInfo.GetComments() {
// 		commentsIds = append(commentsIds, strconv.Itoa(int(comment.GetId())))
// 	}

// 	var likedUserIds []string
// 	for _, like := range postDetailInfo.GetLikes() {
// 		likedUserIds = append(likedUserIds, strconv.Itoa(int(like.GetUserId())))
// 	}

// 	return map[string]interface{}{
// 		"id":                 postDetailInfo.GetId(),
// 		"user_id":            postDetailInfo.GetUserId(),
// 		"content_text":       postDetailInfo.GetContentText(),
// 		"content_image_path": strings.Join(postDetailInfo.GetContentImagePath(), " "),
// 		"visible":            postDetailInfo.GetVisible(),
// 		"created_at":         postDetailInfo.GetCreatedAt(),
// 		"comments_ids":       strings.Join(commentsIds, " "),
// 		"liked_users_ids":    strings.Join(likedUserIds, " "),
// 	}
// }

// func (svc *NewsfeedGenerationService) newRedisLike() map[string]interface{} {
// 	panic("implement me")
// }

// func (svc *NewsfeedGenerationService) newRedisComment() map[string]interface{} {
// 	panic("implement me")
// }
