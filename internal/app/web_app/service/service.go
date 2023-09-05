package service

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/pkg/types"
	"github.com/maxuanquang/social-network/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"
	client_nf "github.com/maxuanquang/social-network/pkg/client/newsfeed"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

var validate = types.NewValidator()

type WebService struct {
	authenticateAndPostClient pb_aap.AuthenticateAndPostClient
	newsfeedClient            pb_nf.NewsfeedClient
	redisClient               *redis.Client

	logger          *zap.Logger
	latencyReporter *prometheus.SummaryVec
	countReporter   *prometheus.CounterVec
}

func NewWebService(cfg *configs.WebConfig) (*WebService, error) {
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	nfClient, err := client_nf.NewClient(cfg.Newsfeed.Hosts)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})
	if redisClient == nil {
		return nil, errors.New("redis connection failed")
	}

	logger, err := utils.NewLogger(&cfg.Logger)
	if err != nil {
		return nil, err
	}

	latencyExporter := promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "social_network_be",
			Subsystem:  "webapp",
			Name:       "latency",
			Help:       "recall latency in milliseconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"component", "status"},
	)

	countExporter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "social_network_be",
			Subsystem: "webapp",
			Name:      "count",
			Help:      "recall count",
		},
		[]string{"component", "type"},
	)

	return &WebService{
		authenticateAndPostClient: aapClient,
		newsfeedClient:            nfClient,
		redisClient:               redisClient,
		logger:                    logger,
		latencyReporter:           latencyExporter,
		countReporter:             countExporter,
	}, nil
}
