package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func MustNew(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	return conn
}
