package authen_and_post

import (
	"context"
	"log"

	"math/rand"

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

// Group: Users

func (a *randomClient) CheckUserAuthentication(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CheckUserAuthentication(ctx, in, opts...)
}

func (a *randomClient) CreateUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CreateUser(ctx, in, opts...)
}

func (a *randomClient) EditUser(ctx context.Context, in *pb.UserDetailInfo, opts ...grpc.CallOption) (*pb.UserResult, error) {
	return a.clients[rand.Intn(len(a.clients))].EditUser(ctx, in, opts...)
}

// Group: Friends

func (a *randomClient) GetUserFollower(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserFollowerInfo, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollower(ctx, in, opts...)
}

func (a *randomClient) FollowUser(ctx context.Context, in *pb.UserAndFollowerInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].FollowUser(ctx, in, opts...)
}

func (a *randomClient) UnfollowUser(ctx context.Context, in *pb.UserAndFollowerInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].UnfollowUser(ctx, in, opts...)
}

func (a *randomClient) GetUserPost(ctx context.Context, in *pb.UserInfo, opts ...grpc.CallOption) (*pb.UserPostDetailInfo, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserPost(ctx, in, opts...)
}

// Group: Posts

func (a *randomClient) CreatePost(ctx context.Context, in *pb.PostDetailInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CreatePost(ctx, in, opts...)
}

func (a *randomClient) GetPost(ctx context.Context, in *pb.PostInfo, opts ...grpc.CallOption) (*pb.PostDetailInfo, error) {
	return a.clients[rand.Intn(len(a.clients))].GetPost(ctx, in, opts...)
}

func (a *randomClient) EditPost(ctx context.Context, in *pb.PostDetailInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].EditPost(ctx, in, opts...)
}

func (a *randomClient) DeletePost(ctx context.Context, in *pb.PostInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].DeletePost(ctx, in, opts...)
}

func (a *randomClient) CommentPost(ctx context.Context, in *pb.CommentInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].CommentPost(ctx, in, opts...)
}

func (a *randomClient) LikePost(ctx context.Context, in *pb.LikeInfo, opts ...grpc.CallOption) (*pb.ActionResult, error) {
	return a.clients[rand.Intn(len(a.clients))].LikePost(ctx, in, opts...)
}
