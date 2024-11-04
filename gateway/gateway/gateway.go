package gateway

import (
	"context"
	pb "github.com/longln/common/api"
)
// this gateway layer manages all connections to other microservices

type OrderGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}