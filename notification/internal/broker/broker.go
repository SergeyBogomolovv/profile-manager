package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service interface{}

type broker struct {
	loginQ    string
	registerQ string
	svc       Service
	logger    *slog.Logger
	ch        *amqp.Channel
}

func MustNew(logger *slog.Logger, conn *amqp.Connection, svc Service) *broker {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	if err := ch.ExchangeDeclare(events.UserExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	loginQ, err := ch.QueueDeclare(events.NotificationLoginQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare login queue: %v", err)
	}

	registerQ, err := ch.QueueDeclare(events.NotificationRegisterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare register queue: %v", err)
	}

	if err := ch.QueueBind(loginQ.Name, events.LoginTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind login queue: %v", err)
	}

	if err := ch.QueueBind(registerQ.Name, events.RegisterTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind register queue: %v", err)
	}

	return &broker{ch: ch, loginQ: loginQ.Name, registerQ: registerQ.Name, svc: svc, logger: logger}
}

func (b *broker) Close() error {
	return b.ch.Close()
}

// Non blocking operation
func (b *broker) Consume(ctx context.Context) {
	go b.consumeLogin(ctx)
	go b.consumeRegister(ctx)
}

func (b *broker) consumeLogin(ctx context.Context) {
	msgs, err := b.ch.Consume(b.loginQ, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume login queue: %v", err)
	}

	for msg := range msgs {
		select {
		case <-ctx.Done():
			return
		default:
			go func() {
				var data events.UserLogin
				if err := json.Unmarshal(msg.Body, &data); err != nil {
					msg.Nack(false, false)
					return
				}
				fmt.Printf("user logined, sending notification, %+v\n", data)
				msg.Ack(false)
			}()
		}
	}
}

func (b *broker) consumeRegister(ctx context.Context) {
	msgs, err := b.ch.Consume(b.registerQ, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume register queue: %v", err)
	}

	for msg := range msgs {
		select {
		case <-ctx.Done():
			return
		default:
			go func() {
				var data events.UserRegister
				if err := json.Unmarshal(msg.Body, &data); err != nil {
					msg.Nack(false, false)
					return
				}
				fmt.Printf("user registered, sending notification, %+v\n", data)
				msg.Ack(false)
			}()
		}
	}
}
