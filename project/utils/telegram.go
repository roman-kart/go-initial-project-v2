package utils

import (
	"time"

	"go.uber.org/zap"
	"gopkg.in/telebot.v3"

	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

type TelegramConfig struct {
	LongPoller struct {
		Timeout uint
	}
}

// TelegramBot provides functionality for creating a [telebot.Bot].
type TelegramBot struct {
	Config              *TelegramConfig
	RabbitMQ            *RabbitMQ
	logger              *zap.Logger
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewTelegram creates a new [TelegramBot].
func NewTelegram(
	config *TelegramConfig,
	logger *zap.Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *TelegramBot {
	return &TelegramBot{
		Config:              config,
		logger:              logger.Named("TelegramBot"),
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("Telegram"),
	}
}

// CreateBot creates a new [telebot.Bot].
func (t *TelegramBot) CreateBot(token string) (*telebot.Bot, error) {
	ew := t.ErrorWrapperCreator.GetMethodWrapper("CreateBot")
	logger := t.logger.Named("Start")

	settings := telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: time.Duration(t.Config.LongPoller.Timeout) * time.Second,
		},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		logger.Error("Failed to create bot", zap.Error(err))
		return nil, ew(err)
	}

	return bot, nil
}
