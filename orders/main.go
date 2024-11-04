package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/longln/common"
	"github.com/longln/common/discovery"
	"github.com/longln/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddress = common.EnvString("GRPC_ADDRESS", "localhost:3000")
	consulAddr = common.EnvString("CONSUL_ADDRESS", "localhost:8500")

)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	err = registry.Register(ctx, instanceID, serviceName, grpcAddress)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := registry.HealthCheck(instanceID, serviceName)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	defer registry.DeRegister(ctx, instanceID, serviceName)
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
	NewGRPCHandler(grpcServer, service)

	if err := grpcServer.Serve(l); err != nil {

	}
}