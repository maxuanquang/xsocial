package service

import (
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	"go.uber.org/zap"

	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"
	// client_nf "github.com/maxuanquang/social-network/pkg/client/newsfeed"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	// pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

var validate = types.NewValidator()

type WebService struct {
	AuthenticateAndPostClient pb_aap.AuthenticateAndPostClient
	// NewsfeedClient            pb_nf.NewsfeedClient
	RedisClient               *redis.Client

	Logger *zap.Logger
}

func NewWebService(conf *configs.WebConfig) (*WebService, error) {
	aapClient, err := client_aap.NewClient(conf.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// nfClient, err := client_nf.NewClient(conf.Newsfeed.Hosts)
	// if err != nil {
	// 	return nil, err
	// }

	redisClient := redis.NewClient(&redis.Options{Addr: conf.Redis.Addr, Password: conf.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	logger, err := newLogger()
	if err != nil {
		return nil, err
	}

	return &WebService{
		AuthenticateAndPostClient: aapClient,
		// NewsfeedClient:            nfClient,
		RedisClient:               redisClient,
		Logger:                    logger,
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
