package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/newsfeed_svc"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "config.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetNewsfeedConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new newsfeed generation service
	gen_service, err := newsfeed_svc.NewNewsfeedGenerationService(cfg)
	go gen_service.Run()

	// Start new newsfeed service
	service, err := newsfeed_svc.NewNewsfeedService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb_nf.RegisterNewsfeedServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}