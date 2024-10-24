package main

import (
	"context"
	"log"
	"net"

	"github.com/longln/common"
	"google.golang.org/grpc"
)

var (
	grpcAddress = common.EnvString("GRPC_ADDRESS", "localhost:3000")
)

func main() {
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)
	service.CreateOrder(context.Background())

	log.Println("GRPC server started at", grpcAddress)
	NewGRPCHandler(grpcServer)

	if err := grpcServer.Serve(l); err != nil {

	}
}