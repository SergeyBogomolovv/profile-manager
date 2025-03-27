package telegram

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	tele "gopkg.in/telebot.v4"
)

type Service interface {
	VerifyTelegram(ctx context.Context, token string, telegramID int64) error
}

type loginer struct {
	logger  *slog.Logger
	bot     *tele.Bot
	service Service
}

func NewLoginer(logger *slog.Logger, bot *tele.Bot, service Service) *loginer {
	return &loginer{logger: logger, bot: bot, service: service}
}

func (l *loginer) Init() {
	l.bot.Handle("/start", l.handleStart)
	l.bot.Handle("/verify", l.handleVerify)
	l.bot.Handle(tele.OnText, l.handleMessage)
}

func (l *loginer) handleStart(c tele.Context) error {
	return c.Send("Привет! Используйте /verify <токен> для подтверждения аккаунта.")
}

func (l *loginer) handleMessage(c tele.Context) error {
	return c.Send("Команда не распознана. Используйте /verify <токен> для подтверждения аккаунта.")
}

func (l *loginer) handleVerify(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Ошибка: укажите токен. Пример: /verify your_token_here")
	}
	token := args[0]

	err := l.service.VerifyTelegram(logger.Inject(context.Background(), l.logger), token, c.Sender().ID)
	if errors.Is(err, domain.ErrInvalidToken) {
		return c.Send("Неверный токен.")
	}
	if errors.Is(err, domain.ErrAccountAlreadyExists) {
		return c.Send("Вы можете привязать только 1 аккаунт.")
	}
	if err != nil {
		l.logger.Error("failed to verify telegram", "error", err)
		return c.Send("Произошла непредвиденная ошибка.")
	}
	return c.Send("Аккаунт успешно привязан!")
}
