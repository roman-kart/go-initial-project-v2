package managers

import (
	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/utils"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type TelegramBotManager struct {
	Config             *config.Config
	Logger             *utils.Logger
	logger             *zap.Logger
	TelegramBot        *utils.TelegramBot
	telegramBot        *telebot.Bot
	StatManager        *StatManager
	UserAccountManager *UserAccountManager
}

func NewTelegramBotManager(
	config *config.Config,
	logger *utils.Logger,
	statManager *StatManager,
	userAccountManager *UserAccountManager,
	telegramBot *utils.TelegramBot,
) *TelegramBotManager {
	return &TelegramBotManager{
		Config:             config,
		Logger:             logger,
		logger:             logger.Logger.Named("TelegramBotManager"),
		TelegramBot:        telegramBot,
		StatManager:        statManager,
		UserAccountManager: userAccountManager,
	}
}

func (t *TelegramBotManager) Prepare() error {
	return t.createBot()
}

func (t *TelegramBotManager) createBot() error {
	bot, err := t.TelegramBot.CreateBotDefault()
	if err != nil {
		return err
	}
	t.telegramBot = bot
	return nil
}

func (t *TelegramBotManager) GetBot() *telebot.Bot {
	return t.telegramBot
}
