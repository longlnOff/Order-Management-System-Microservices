package main


import (
	"context"
	pb "github.com/longln/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}