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
func (a *randomClient) CheckUserAuthentication(ctx context.Context, in *pb.CheckUserAuthenticationRequest, opts ...grpc.CallOption) (*pb.CheckUserAuthenticationResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CheckUserAuthentication(ctx, in, opts...)
}

func (a *randomClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest, opts ...grpc.CallOption) (*pb.CreateUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CreateUser(ctx, in, opts...)
}

func (a *randomClient) EditUser(ctx context.Context, in *pb.EditUserRequest, opts ...grpc.CallOption) (*pb.EditUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditUser(ctx, in, opts...)
}

func (a *randomClient) GetUserDetailInfo(ctx context.Context, in *pb.GetUserDetailInfoRequest, opts ...grpc.CallOption) (*pb.GetUserDetailInfoResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserDetailInfo(ctx, in, opts...)
}

// Group: Friends

func (a *randomClient) GetUserFollower(ctx context.Context, in *pb.GetUserFollowerRequest, opts ...grpc.CallOption) (*pb.GetUserFollowerResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollower(ctx, in, opts...)
}

func (a *randomClient) GetUserFollowing(ctx context.Context, in *pb.GetUserFollowingRequest, opts ...grpc.CallOption) (*pb.GetUserFollowingResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserFollowing(ctx, in, opts...)
}

func (a *randomClient) FollowUser(ctx context.Context, in *pb.FollowUserRequest, opts ...grpc.CallOption) (*pb.FollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].FollowUser(ctx, in, opts...)
}

func (a *randomClient) UnfollowUser(ctx context.Context, in *pb.UnfollowUserRequest, opts ...grpc.CallOption) (*pb.UnfollowUserResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].UnfollowUser(ctx, in, opts...)
}

func (a *randomClient) GetUserPosts(ctx context.Context, in *pb.GetUserPostsRequest, opts ...grpc.CallOption) (*pb.GetUserPostsResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetUserPosts(ctx, in, opts...)
}

// Group: Posts

func (a *randomClient) CreatePost(ctx context.Context, in *pb.CreatePostRequest, opts ...grpc.CallOption) (*pb.CreatePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CreatePost(ctx, in, opts...)
}

func (a *randomClient) GetPostDetailInfo(ctx context.Context, in *pb.GetPostDetailInfoRequest, opts ...grpc.CallOption) (*pb.GetPostDetailInfoResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetPostDetailInfo(ctx, in, opts...)
}

func (a *randomClient) EditPost(ctx context.Context, in *pb.EditPostRequest, opts ...grpc.CallOption) (*pb.EditPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].EditPost(ctx, in, opts...)
}

func (a *randomClient) DeletePost(ctx context.Context, in *pb.DeletePostRequest, opts ...grpc.CallOption) (*pb.DeletePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].DeletePost(ctx, in, opts...)
}

func (a *randomClient) CommentPost(ctx context.Context, in *pb.CommentPostRequest, opts ...grpc.CallOption) (*pb.CommentPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].CommentPost(ctx, in, opts...)
}

func (a *randomClient) LikePost(ctx context.Context, in *pb.LikePostRequest, opts ...grpc.CallOption) (*pb.LikePostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].LikePost(ctx, in, opts...)
}

func (a *randomClient) GetS3PresignedUrl(ctx context.Context, in *pb.GetS3PresignedUrlRequest, opts ...grpc.CallOption) (*pb.GetS3PresignedUrlResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetS3PresignedUrl(ctx, in, opts...)
}
