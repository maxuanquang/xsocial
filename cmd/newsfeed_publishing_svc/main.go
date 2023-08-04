package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/newsfeed_publishing_svc"
	pb_nfp "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "config.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetNewsfeedPublishingConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new newsfeed publishing service
	service, err := newsfeed_publishing_svc.NewNewsfeedPublishingService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	// Run fanout worker
	go service.Run()

	// Start grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb_nfp.RegisterNewsfeedPublishingServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
