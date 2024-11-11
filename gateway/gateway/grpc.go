package gateway

import (
	"context"
	"log"

	pb "github.com/longln/common/api"
	"github.com/longln/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{
		registry: registry,
	}
}

func (g *gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	order, err := c.CreateOrder(ctx, 
								&pb.CreateOrderRequest{
									CustomerID: req.CustomerID,
									Items:      req.Items,
								})
	
	log.Printf("order: %v", order)
	return order, err
}



func (g *gateway) GetOrder(ctx context.Context, orderID string, customerID string) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	
	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderID: orderID,
		CustomerID: customerID,
	})	
}
