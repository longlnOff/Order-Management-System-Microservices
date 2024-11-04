package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(user string, password string, host string, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	conn, err := amqp.Dial(address)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(OrderCreatedPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	return ch, conn.Close
}