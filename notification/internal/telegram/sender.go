package telegram

import (
	"errors"

	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	tele "gopkg.in/telebot.v4"
)

type Sender interface {
	SendLoginNotification(telegramID int64, data domain.LoginNotification) error
}

type sender struct {
	bot *tele.Bot
}

func NewSender(bot *tele.Bot) Sender {
	return &sender{bot: bot}
}

func (s *sender) SendLoginNotification(telegramID int64, data domain.LoginNotification) error {
	_, err := s.bot.Send(tele.ChatID(telegramID), "Login\n\n"+"IP: "+data.IP+"\nTime: "+data.Time+"\nType: "+data.Type)
	if errors.Is(err, tele.ErrBlockedByUser) {
		return nil
	}
	return err
}
