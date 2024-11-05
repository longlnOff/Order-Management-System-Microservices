package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/longln/common/api"
	"github.com/longln/common/broker"
	"github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{
		service: service,
	}
}

func (c *consumer) Listen (ch *amqp091.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent,
		false,
		false,
		false,
		false,
			nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// var forever chan struct{}

	go func() {
		for msg := range msgs {
			log.Printf("received a message: %v", msg.Body)
			o := &pb.Order{}
			if err := json.Unmarshal(msg.Body, o); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)
				continue
			}
			log.Printf("Payment link: %s", paymentLink)
			// err = ch.PublishWithContext(
			// 	context.Background(),
			// 	"",
			// 	broker.OrderCreatedPaid,
			// 	false,
			// 	false,
			// 	amqp091.Publishing{
			// 		ContentType: "text/plain",
			// 		Body:        []byte(paymentLink),
			// 		DeliveryMode: amqp091.Persistent,
			// 	})
			// if err != nil {
			// 	log.Printf("failed to publish payment link: %v", err)
			// 	continue
			// }

		}
	}()


}	