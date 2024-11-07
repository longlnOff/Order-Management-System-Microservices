package main

import (
	"context"
	"testing"

	pb "github.com/longln/common/api"
	"github.com/longln/payments/processor/inmem"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	service := NewService(processor)

	t.Run("should create payment link", func(t *testing.T) {
		link, err := service.CreatePayment(context.Background(), &pb.Order{})
		require.NoError(t, err)
		require.NotEmpty(t, link)
	})
}