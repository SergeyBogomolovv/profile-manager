package telegram

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	tele "gopkg.in/telebot.v4"
)

type SetupService interface {
	LinkTelegram(ctx context.Context, token string, telegramID int64) error
	UnlinkTelegram(ctx context.Context, telegramID int64) error
	UpdateSubscriptionStatus(ctx context.Context, telegramID int64, subType domain.SubscriptionType, enabled bool) error
}

type loginer struct {
	logger  *slog.Logger
	bot     *tele.Bot
	service SetupService
	state   *state
}

func NewLoginer(logger *slog.Logger, bot *tele.Bot, service SetupService) *loginer {
	return &loginer{logger: logger, bot: bot, service: service, state: NewState()}
}

func (l *loginer) Init() {
	l.bot.Handle("/start", l.handleStart)
	l.bot.Handle("/link", l.startLinkTg)
	l.bot.Handle("/unlink", l.handleUnlinkTelegram)
	l.bot.Handle("/cancel", l.handleCancel)
	l.bot.Handle("/enable", l.startEnableNotifications)
	l.bot.Handle("/disable", l.startDisableNotifications)
	l.bot.Handle(tele.OnText, l.handleMessage)
}

func (l *loginer) handleStart(c tele.Context) error {
	return c.Send("Привет! Используйте /link для привязки аккаунта.")
}

func (l *loginer) handleMessage(c tele.Context) error {
	state, _ := l.state.Get(c.Sender().ID)

	switch state {
	case stateWaitingToken:
		return l.handleLinkTelegram(c)
	case stateWaitingSubTypeEnable:
		return l.handleEnableNotifications(c)
	case stateWaitingSubTypeDisable:
		return l.handleDisableNotifications(c)
	default:
		return c.Send("Команда не распознана.")
	}
}

func (l *loginer) handleCancel(c tele.Context) error {
	_, ok := l.state.Get(c.Sender().ID)
	if !ok {
		return c.Send("Активные действия не найдены.")
	}
	return l.sendAndClear(c, "Действие отменено.")
}

func (l *loginer) startLinkTg(c tele.Context) error {
	l.state.Set(c.Sender().ID, stateWaitingToken)
	return c.Send("Введите токен. Для отмены используйте /cancel")
}

func (l *loginer) handleLinkTelegram(c tele.Context) error {
	token := c.Message().Text
	userID := c.Sender().ID
	err := l.service.LinkTelegram(l.loggerCtx(), token, userID)
	if errors.Is(err, domain.ErrInvalidToken) {
		return c.Send("Неверный токен.")
	}
	if errors.Is(err, domain.ErrAccountAlreadyExists) {
		return l.sendAndClear(c, "Вы можете привязать только 1 аккаунт.")
	}
	if errors.Is(err, domain.ErrActionDontNeeded) {
		return l.sendAndClear(c, "Аккаунт уже привязан.")
	}
	if err != nil {
		l.logger.Error("failed to link telegram", "error", err)
		return c.Send("Произошла непредвиденная ошибка.")
	}
	return l.sendAndClear(c, "Аккаунт успешно привязан!")
}

func (l *loginer) handleUnlinkTelegram(c tele.Context) error {
	err := l.service.UnlinkTelegram(l.loggerCtx(), c.Sender().ID)
	if errors.Is(err, domain.ErrActionDontNeeded) {
		return c.Send("Аккаунт не привязан.")
	}
	if err != nil {
		l.logger.Error("failed to unlink telegram", "error", err)
		return c.Send("Произошла непредвиденная ошибка.")
	}
	return c.Send("Аккаунт успешно отвязан!")
}

func (l *loginer) startEnableNotifications(c tele.Context) error {
	l.state.Set(c.Sender().ID, stateWaitingSubTypeEnable)
	return c.Send("Выберите тип уведомлений. Для отмены используйте /cancel", accountTypeMenu)
}

func (l *loginer) handleEnableNotifications(c tele.Context) error {
	userID := c.Sender().ID
	subType, err := getSubscriptionType(c.Message().Text)
	if err != nil {
		return c.Send("Выберите тип уведомлений: Телеграм или Почта.")
	}

	err = l.service.UpdateSubscriptionStatus(l.loggerCtx(), userID, subType, true)
	if errors.Is(err, domain.ErrActionDontNeeded) {
		return l.sendAndClear(c, "Уведомления уже подключены.")
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return l.sendAndClear(c, "Привяжите аккаунт с помощью команды /link")
	}
	if err != nil {
		l.logger.Error("failed to enable notifications", "error", err)
		return c.Send("Произошла непредвиденная ошибка.")
	}
	return l.sendAndClear(c, "Уведомления успешно подключены.")
}

func (l *loginer) startDisableNotifications(c tele.Context) error {
	l.state.Set(c.Sender().ID, stateWaitingSubTypeDisable)
	return c.Send("Выберите тип уведомлений. Для отмены используйте /cancel", accountTypeMenu)
}

func (l *loginer) handleDisableNotifications(c tele.Context) error {
	userID := c.Sender().ID
	subType, err := getSubscriptionType(c.Message().Text)
	if err != nil {
		return c.Send("Выберите тип уведомлений: Телеграм или Почта.")
	}
	err = l.service.UpdateSubscriptionStatus(l.loggerCtx(), userID, subType, false)
	if errors.Is(err, domain.ErrActionDontNeeded) {
		return l.sendAndClear(c, "Уведомления уже отключены.")
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return l.sendAndClear(c, "Привяжите аккаунт с помощью команды /link")
	}
	if err != nil {
		l.logger.Error("failed to disable notifications", "error", err)
		return c.Send("Произошла непредвиденная ошибка.")
	}

	return l.sendAndClear(c, "Уведомления успешно отключены.")
}

func getSubscriptionType(text string) (domain.SubscriptionType, error) {
	switch text {
	case AccountTypeTelegram:
		return domain.SubscriptionTypeTelegram, nil
	case AccountTypeEmail:
		return domain.SubscriptionTypeEmail, nil
	default:
		return "", errors.New("unknown subscription type")
	}
}

func (l *loginer) loggerCtx() context.Context {
	return logger.Inject(context.Background(), l.logger)
}

func (l *loginer) sendAndClear(c tele.Context, message string) error {
	l.state.Delete(c.Sender().ID)
	return c.Send(message, clearMenu)
}
