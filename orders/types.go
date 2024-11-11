package main

import (
	"context"
	pb "github.com/longln/common/api"

)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
}


type OrderStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, id string, customerID string) (*pb.Order, error)
}