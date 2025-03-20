package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type broker struct {
	ch *amqp.Channel
}

func MustNew(conn *amqp.Connection) *broker {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	if err := ch.ExchangeDeclare(events.UserExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}
	return &broker{ch: ch}
}

func (b *broker) Close() error {
	return b.ch.Close()
}

func (b *broker) PublishUserRegister(user events.UserRegister) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	return b.ch.Publish(events.UserExchange, events.RegisterTopic, false, false, msg)
}

func (b *broker) PublishUserLogin(user events.UserLogin) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	return b.ch.Publish(events.UserExchange, events.LoginTopic, false, false, msg)
}
