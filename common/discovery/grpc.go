package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, servicename string, registry Registry) (*grpc.ClientConn, error) {
	addrss, err := registry.Discover(ctx, servicename)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(addrss[rand.Intn(len(addrss))], grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}