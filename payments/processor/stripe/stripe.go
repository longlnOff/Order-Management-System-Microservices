package stripe

import (
	pb "github.com/longln/common/api"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81"

)
type Stripe struct {

}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(order *pb.Order) (string, error) {

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
						Price: stripe.String(item.PriceID),
						Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	gatewaySuccessURL := ""
	params := &stripe.CheckoutSessionParams{
		LineItems: items,
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
