package processor

import (
	pb "github.com/longln/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}