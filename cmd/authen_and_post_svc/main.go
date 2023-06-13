package main

import (
	"fmt"
	"log"
	"net"

	"github.com/maxuanquang/social-network/internal/app/authen_and_post_svc"
	"github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"google.golang.org/grpc"
)

func main() {
	// Start authenticate and post service
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 1080))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	service := authen_and_post_svc.NewAuthenticateAndPostService()
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...) // The `...` is used to unpack the slice into individual arguments
	authen_and_post.RegisterAuthenticateAndPostServer(grpcServer, service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
