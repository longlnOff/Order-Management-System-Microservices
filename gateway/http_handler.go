package main

import (
	"errors"
	"net/http"

	"github.com/longln/common"
	pb "github.com/longln/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	// gateway
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{
		client: client,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}


func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}


	order, err := h.client.CreateOrder(r.Context(), 
										&pb.CreateOrderRequest{
											CustomerID: customerID,
											Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusInternalServerError, rStatus.Message())
			return
		}

	}

	common.WriteJSON(w, http.StatusCreated, order)

}

func validateItems(items [] *pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items must not be empty")
	}

	for  _, i := range items {
		if i.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}
		if i.ItemID == "" {
			return errors.New("itemID must not be empty")
		}
	}

	return nil
}