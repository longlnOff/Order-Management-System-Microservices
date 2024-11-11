package main

import (
	"context"

	"github.com/longln/common"
	pb "github.com/longln/common/api"
)

type service struct {
	store OrderStore
}


func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) GetOrder(ctx context.Context, r *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.store.Get(ctx, r.OrderID, r.CustomerID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {

	id, err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	order := &pb.Order{
		ID: id,
		CustomerID: p.CustomerID,
		Status: "pending",
		Items: items,
	}
	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	merge := mergeItemsQuantities(p.Items)

	// validate with the stock service

	// Temporary
	var itemsWithPrice []*pb.Item
	for _, i := range(merge) {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			ID: i.ItemID,
			PriceID: "price_1QIS1pFfewWNTBtoJLn9m0qh",
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
}


func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, i := range items {
		found := false
		for _, j := range merged {
			if i.ItemID == j.ItemID {
				j.Quantity += i.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, i)
		}
	}

	return merged
}
