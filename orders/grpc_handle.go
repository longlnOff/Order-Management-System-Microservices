package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/longln/common/api"
	"github.com/longln/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := h.service.GetOrder(ctx, r)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	// order := &pb.Order{ID: "42", Items: r.Items}
	items, err := h.service.ValidateOrder(ctx, r)
	if err != nil {
		return nil, err
	}
	order, err := h.service.CreateOrder(ctx, r, items)
	
	log.Printf("New order received!\nOrder: %v\n", order)

	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}
	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	err = h.channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledOrder,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return nil, err
	}


	return order, nil
}