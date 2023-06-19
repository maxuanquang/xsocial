package authen_and_post

import (
	"context"
	"log"

	"math/rand"

	// "github.com/maxuanquang/social-network/configs"
	pb "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(hosts []string) (pb.AuthenticateAndPostClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]pb.AuthenticateAndPostClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
			return nil, err
		}

		client := pb.NewAuthenticateAndPostClient(conn)
		clients = append(clients, client)
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb.AuthenticateAndPostClient
}

func (a *randomClient) CheckUserAuthentication(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CheckUserAuthentication(ctx, in, opts...)
}

func (a *randomClient) CreateUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CreateUser(ctx, in, opts...)
}

func (a *randomClient) EditUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].EditUser(ctx, in, opts...)
}

func (a *randomClient) GetUserFollower(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserFollower, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollower(ctx, in, opts...)
}

func (a *randomClient) GetPostDetail(ctx context.Context, in *pb.GetPostRequest, opts ...grpc.CallOption) (*pb.Post, error) {
	return a.clients[rand.Intn(len(a.clients))].GetPostDetail(ctx, in, opts...)
}
