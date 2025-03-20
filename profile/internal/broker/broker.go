package broker

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ProfileService interface {
	Create(ctx context.Context, user events.UserRegister) error
}

type broker struct {
	qName   string
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

	q, err := ch.QueueDeclare(events.ProfileRegisterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	if err := ch.QueueBind(q.Name, events.RegisterTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	return &broker{ch: ch, qName: q.Name, profile: profile, logger: logger}
}

func (b *broker) Close() error {
	return b.ch.Close()
}

func (b *broker) Consume(ctx context.Context) {
	msgs, err := b.ch.Consume(b.qName, "profile-service", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume queue: %v", err)
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
				if err := b.profile.Create(ctx, data); err != nil {
					b.logger.Error("failed to create profile", "error", err)
					msg.Nack(false, false)
					return
				}
				msg.Ack(false)
			}()
		}
	}
}
