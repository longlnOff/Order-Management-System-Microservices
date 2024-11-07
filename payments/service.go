package main

import (
	"context"
	// "log"

	pb "github.com/longln/common/api"
	"github.com/longln/payments/processor"
)
type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{
		processor: processor,
	}
}

func (s *service) CreatePayment(ctx context.Context,order *pb.Order) (string, error) {
	// connect to payment processor
	link, err := s.processor.CreatePaymentLink(order)
	if err != nil {
		return "", err
	}
	// log.Printf("payment link: %s", link)
	return link, nil
}
