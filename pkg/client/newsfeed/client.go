package newsfeed

import (
	"context"
	"log"
	"math/rand"

	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(hosts []string) (pb_nf.NewsfeedClient, error) {
	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	clients := make([]pb_nf.NewsfeedClient, 0, len(hosts))
	for _, host := range hosts {
		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
			return nil, err
		}

		client := pb_nf.NewNewsfeedClient(conn)
		clients = append(clients, client)
	}

	return &randomClient{clients: clients}, nil
}

type randomClient struct {
	clients []pb_nf.NewsfeedClient
}

func (a *randomClient) GetNewsfeed(ctx context.Context, in *pb_nf.NewsfeedRequest, opts ...grpc.CallOption) (*pb_nf.NewsfeedResponse, error) {
	return a.clients[rand.Intn(len(a.clients))].GetNewsfeed(ctx, in, opts...)
}
