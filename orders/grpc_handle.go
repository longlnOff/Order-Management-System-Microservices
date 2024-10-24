package main

import (
	"context"
	// "fmt"
	"log"

	pb "github.com/longln/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGRPCHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received!\nOrder: %v\n", r)
	// return nil, fmt.Errorf("not implemented")
	return &pb.Order{ID: "42"}, nil
}