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

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}
	order := &pb.Order{
		ID: "42",
		CustomerID: p.CustomerID,
		Items: items,
		Status: "pending",
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
