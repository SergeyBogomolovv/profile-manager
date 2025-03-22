package broker

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProfileService interface {
	Create(ctx context.Context, user events.UserRegister) error
}

type broker struct {
	profile ProfileService
	logger  *slog.Logger
	ch      *amqp.Channel
}

func MustNew(logger *slog.Logger, conn *amqp.Connection, profile ProfileService) *broker {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	if err := ch.ExchangeDeclare(events.UserExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	return &broker{ch: ch, profile: profile, logger: logger}
}

func (b *broker) Close() error {
	return b.ch.Close()
}

// Non blocking operation
func (b *broker) Consume(ctx context.Context) {
	go b.consumeRegister(ctx)
}

func (b *broker) consumeRegister(ctx context.Context) {
	q, err := b.ch.QueueDeclare(events.ProfileRegisterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	if err := b.ch.QueueBind(q.Name, events.RegisterTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := b.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume queue: %v", err)
	}

	for msg := range msgs {
		select {
		case <-ctx.Done():
			return
		default:
			go b.handleRegister(logger.Inject(ctx, b.logger), msg)
		}
	}
}

func (b *broker) handleRegister(ctx context.Context, msg amqp.Delivery) {
	var data events.UserRegister
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		msg.Nack(false, true)
		return
	}
	if err := b.profile.Create(ctx, data); err != nil {
		b.logger.Error("failed to create profile", "error", err)
		msg.Nack(false, true)
		return
	}
	msg.Ack(false)
}
