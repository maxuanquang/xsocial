package authen_and_post

import (
	"context"
	"log"

	"math/rand"

	pb "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed_publishing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(hosts []string) (pb.NewsfeedPublishingClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]pb.NewsfeedPublishingClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
			return nil, err
		}

		client := pb.NewNewsfeedPublishingClient(conn)
		clients = append(clients, client)
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb.NewsfeedPublishingClient
}

func (a *randomClient) PublishPost(ctx context.Context, in *pb.PublishPostRequest, opts ...grpc.CallOption) (*pb.PublishPostResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].PublishPost(ctx, in, opts...)
}