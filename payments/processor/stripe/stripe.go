package stripe

import (
	"fmt"
	"log"

	"github.com/longln/common"
	pb "github.com/longln/common/api"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)


var gatewayHTTPAddress = common.EnvString("GATEWAY_HTTP_ADDRESS", "http://localhost:8080")

type Stripe struct {

}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(order *pb.Order) (string, error) {
	log.Printf("create payment link for order: %v", order)
	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
						Price: stripe.String("price_1QIS1pFfewWNTBtoJLn9m0qh"),
						Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	gatewaySuccessURL := fmt.Sprintf("%s/success.html", gatewayHTTPAddress)
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
