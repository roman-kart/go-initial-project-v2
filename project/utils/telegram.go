package utils

import (
	"time"

	"go.uber.org/zap"
	"gopkg.in/telebot.v3"

	c "github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
)

// TelegramBot provides functionality for creating a [telebot.Bot].
type TelegramBot struct {
	Config              *c.Config
	Logger              *Logger
	RabbitMQ            *RabbitMQ
	logger              *zap.Logger
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewTelegram creates a new [TelegramBot].
func NewTelegram(
	config *c.Config,
	logger *Logger,
	rabbitMQ *RabbitMQ,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *TelegramBot {
	return &TelegramBot{
		Config:              config,
		Logger:              logger,
		RabbitMQ:            rabbitMQ,
		logger:              logger.Logger.Named("TelegramBot"),
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
			Timeout: time.Duration(t.Config.Telegram.LongPoller.Timeout) * time.Second,
		},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		logger.Error("Failed to create bot", zap.Error(err))
		return nil, ew(err)
	}

	return bot, nil
}

// CreateBotDefault creates a new [telebot.Bot] with default token.
func (t *TelegramBot) CreateBotDefault() (*telebot.Bot, error) {
	return t.CreateBot(t.Config.Telegram.Token)
}
