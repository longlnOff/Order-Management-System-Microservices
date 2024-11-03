package main

import (
	"context"
	"log"

	"github.com/longln/common"
	pb "github.com/longln/common/api"
)

type service struct {
	store OrderStore
}


func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}

	mergeditems := mergeItemsQuantities(p.Items)
	log.Println(mergeditems)

	// validate with the stock service

	return nil
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
