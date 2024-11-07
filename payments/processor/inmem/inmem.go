package inmem

import (
	pb "github.com/longln/common/api"
)

type Inmem struct {

}

func NewInmem() *Inmem {
	return &Inmem{}
}

func (i *Inmem) CreatePaymentLink(order *pb.Order) (string, error) {
	return "dummy-link", nil
}