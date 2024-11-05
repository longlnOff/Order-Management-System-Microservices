package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/longln/common"
	"github.com/longln/common/broker"
	"github.com/longln/common/discovery"
	"github.com/longln/common/discovery/consul"
	"google.golang.org/grpc"
	"github.com/stripe/stripe-go/v81"
)

var (
	serviceName = "payments"
	grpcAddress = common.EnvString("GRPC_ADDRESS", "localhost:3002")
	consulAddr = common.EnvString("CONSUL_ADDRESS", "localhost:8500")
	amqpUser = common.EnvString("AMQP_USER", "guest")
	amqpPassword = common.EnvString("AMQP_PASSWORD", "guest")
	amqpHost = common.EnvString("AMQP_HOST", "localhost")
	amqpPort = common.EnvString("AMQP_PORT", "5672")
	stripeKey = common.EnvString("STRIPE_KEY", "")
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
	// stripe setup
	stripe.Key = stripeKey

	// broker connection
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	svc := NewService()

	amqpConsummer := NewConsumer(svc)

	go amqpConsummer.Listen(ch)

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer l.Close()

	log.Println("GRPC server started at", grpcAddress)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("failed to serve", err.Error())
	}
}