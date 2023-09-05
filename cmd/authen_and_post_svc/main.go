package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/authen_and_post_svc"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"google.golang.org/grpc"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "config.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetAuthenticateAndPostConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// Start new authenticate and post service
	service, err := authen_and_post_svc.NewAuthenticateAndPostService(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("can not listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb_aap.RegisterAuthenticateAndPostServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
