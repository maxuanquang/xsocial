package newsfeed_svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"github.com/segmentio/kafka-go"
)

type NewsfeedGenerationService struct {
	kafkaReader               *kafka.Reader
	redisClient               *redis.Client
	authenticateAndPostClient pb_aap.AuthenticateAndPostClient
}

func NewNewsfeedGenerationService(cfg *configs.NewsfeedConfig) (*NewsfeedGenerationService, error) {
	// Connect to kafka
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		Logger:  log.New(os.Stdout, "kafka reader: ", 0),
	})
	if kafkaReader == nil {
		return nil, errors.New("kafka connection failed")
	}

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	// Connect to aapClient
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// Return
	return &NewsfeedGenerationService{
		kafkaReader:               kafkaReader,
		redisClient:               redisClient,
		authenticateAndPostClient: aapClient,
	}, nil
}

func (svc *NewsfeedGenerationService) Run() {
	for {
		// Take message from kafka
		message, err := svc.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Process message
		svc.processMessage(message)
	}
}

func (svc *NewsfeedGenerationService) processMessage(message kafka.Message) {
	jsonData := message.Value
	msgType := string(message.Key)

	// Process message based on its key
	if msgType == "post" {
		svc.processPost(jsonData)
	}
}

func (svc *NewsfeedGenerationService) processPost(jsonData []byte) {
	var postDetailInfo pb_aap.PostDetailInfo
	err := json.Unmarshal(jsonData, &postDetailInfo)
	if err != nil {
		fmt.Println(err)
	}

	// Try to get distributed lock
	lockName := "lock-post-" + strconv.Itoa(int(postDetailInfo.GetId()))
	svc.acquireDistributedLock(lockName)
	defer svc.releaseDistributedLock(lockName)

	// Do following works
	// 1. Create a post object in redis
	postKey := "post-" + strconv.Itoa(int(postDetailInfo.GetId()))
	mapRedisPost := svc.newMapRedisPost(&postDetailInfo)
	_, err = svc.redisClient.HSet(context.Background(), postKey, mapRedisPost).Result()
	if err != nil {
		panic(err)
	}

	// 2. Find followers of user that created post
	followersKey := "followers-" + strconv.Itoa(int(postDetailInfo.UserId))
	numKey, _ := svc.redisClient.Exists(context.Background(), followersKey).Result()
	if numKey == 0 {
		userFollowersInfo, err := svc.authenticateAndPostClient.GetUserFollower(context.Background(), &pb_aap.UserInfo{
			Id: postDetailInfo.GetId(),
		})
		if err != nil {
			panic(err)
		}

		var followersIDs []string
		for _, userInfo := range userFollowersInfo.Followers {
			followersIDs = append(followersIDs, strconv.Itoa(int(userInfo.GetId())))
		}

		if len(followersIDs) > 0 {
			_, err = svc.redisClient.RPush(context.Background(), followersKey, followersIDs).Result()
			if err != nil {
				panic(err)
			}
		}
	}
	followersIDs, err := svc.redisClient.LRange(context.Background(), followersKey, 0, -1).Result()
	if err != nil {
		panic("err")
	}

	// 3. Add this post_id into followers' newsfeed
	for _, id := range followersIDs {
		newsfeedKey := "newsfeed-" + id
		_, err := svc.redisClient.RPush(context.Background(), newsfeedKey, postDetailInfo.GetId()).Result()
		if err != nil {
			panic(err)
		}
	}
}

func (svc *NewsfeedGenerationService) acquireDistributedLock(lockName string) {
	for {
		err := svc.redisClient.SetNX(context.Background(), lockName, "lock", 5*time.Minute).Err()
		if err == nil {
			break
		}
		time.Sleep(time.Second / 2)
		continue
	}
}

func (svc *NewsfeedGenerationService) releaseDistributedLock(lockName string) {
	svc.redisClient.Del(context.Background(), lockName)
}

func (svc *NewsfeedGenerationService) newMapRedisPost(postDetailInfo *pb_aap.PostDetailInfo) map[string]interface{} {
	var commentsIds []string
	for _, comment := range postDetailInfo.GetComments() {
		commentsIds = append(commentsIds, strconv.Itoa(int(comment.GetId())))
	}

	var likedUserIds []string
	for _, like := range postDetailInfo.GetLikes() {
		likedUserIds = append(likedUserIds, strconv.Itoa(int(like.GetUserId())))
	}

	return map[string]interface{}{
		"id":                 postDetailInfo.GetId(),
		"user_id":            postDetailInfo.GetUserId(),
		"content_text":       postDetailInfo.GetContentText(),
		"content_image_path": strings.Join(postDetailInfo.GetContentImagePath(), " "),
		"visible":            postDetailInfo.GetVisible(),
		"create_at":          postDetailInfo.GetCreatedAt(),
		"comments_ids":       strings.Join(commentsIds, " "),
		"liked_user_ids":     strings.Join(likedUserIds, " "),
	}
}

func (svc *NewsfeedGenerationService) newRedisLike() map[string]interface{} {
	panic("implement me")
}

func (svc *NewsfeedGenerationService) newRedisComment() map[string]interface{} {
	panic("implement me")
}
