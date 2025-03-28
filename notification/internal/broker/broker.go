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

type NotifyService interface {
	SendLoginNotification(ctx context.Context, data events.UserLogin) error
	HandleRegister(ctx context.Context, data events.UserRegister) error
}

type broker struct {
	svc    NotifyService
	logger *slog.Logger
	ch     *amqp.Channel
}

func MustNew(logger *slog.Logger, conn *amqp.Connection, svc NotifyService) *broker {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	if err := ch.ExchangeDeclare(events.UserExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	return &broker{ch: ch, svc: svc, logger: logger}
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
	q, err := b.ch.QueueDeclare(events.NotificationLoginQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare login queue: %v", err)
	}

	if err := b.ch.QueueBind(q.Name, events.LoginTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind login queue: %v", err)
	}

	msgs, err := b.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume login queue: %v", err)
	}

	for msg := range msgs {
		select {
		case <-ctx.Done():
			return
		default:
			go b.handleLogin(logger.Inject(ctx, b.logger), msg)
		}
	}
}

func (b *broker) handleLogin(ctx context.Context, msg amqp.Delivery) {
	var data events.UserLogin
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		msg.Nack(false, true)
		return
	}
	if err := b.svc.SendLoginNotification(ctx, data); err != nil {
		logger.Extract(ctx).Error("failed to send login notification", "error", err)
		msg.Nack(false, true)
		return
	}
	msg.Ack(false)
}

func (b *broker) consumeRegister(ctx context.Context) {
	q, err := b.ch.QueueDeclare(events.NotificationRegisterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare register queue: %v", err)
	}

	if err := b.ch.QueueBind(q.Name, events.RegisterTopic, events.UserExchange, false, nil); err != nil {
		log.Fatalf("failed to bind register queue: %v", err)
	}

	msgs, err := b.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to consume register queue: %v", err)
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
	if err := b.svc.HandleRegister(ctx, data); err != nil {
		logger.Extract(ctx).Error("failed to handle register", "error", err)
		msg.Nack(false, true)
		return
	}
	msg.Ack(false)
}
