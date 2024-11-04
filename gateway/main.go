package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/longln/common"
	"github.com/longln/common/discovery"
	"github.com/longln/common/discovery/consul"
	"github.com/longln/gateway/gateway"
)



var (
	serviceName = "gateway"
	httpAddress = common.EnvString("HTTP_ADDRESS", ":8080")
	consulAddr = common.EnvString("CONSUL_ADDRESS", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	err = registry.Register(ctx, instanceID, serviceName, httpAddress)
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
	

	ordersGateway := gateway.NewGRPCGateway(registry)
	mux := http.NewServeMux()
	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("starting server at %s", httpAddress)
	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatal("failed to start server", err)
	}
}