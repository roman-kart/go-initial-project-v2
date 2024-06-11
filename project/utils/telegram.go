package utils

import (
	"time"

	c "github.com/roman-kart/go-initial-project/project/config"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type TelegramBot struct {
	Config   *c.Config
	Logger   *Logger
	RabbitMQ *RabbitMQ
	logger   *zap.Logger
}

// New instance creation

func NewTelegram(
	config *c.Config,
	logger *Logger,
	rabbitMQ *RabbitMQ,
) *TelegramBot {
	return &TelegramBot{
		Config:   config,
		Logger:   logger,
		RabbitMQ: rabbitMQ,
		logger:   logger.Logger.Named("TelegramBot"),
	}
}

func (t *TelegramBot) CreateBot(token string) (*telebot.Bot, error) {
	logger := t.logger.Named("Start")

	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		logger.Error("Failed to create bot", zap.Error(err))
		return nil, err
	}

	return bot, nil
}

func (t *TelegramBot) CreateBotDefault() (*telebot.Bot, error) {
	return t.CreateBot(t.Config.Telegram.Token)
}
