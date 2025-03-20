package app

import (
	"log/slog"

	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
)

type app struct {
	logger *slog.Logger
	conf   *config.Config
}

func New(logger *slog.Logger, conf *config.Config) *app {
	return &app{logger: logger, conf: conf}
}

func (a *app) Start() {
	a.logger.Info("notification service started")
}

func (a *app) Stop() {
	a.logger.Info("notification service stopped")
}
