package managers

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// TelegramBotManager managing [utils.TelegramBot].
type TelegramBotManager struct {
	Config              *config.Config
	Logger              *utils.Logger
	logger              *zap.Logger
	TelegramBot         *utils.TelegramBot
	telegramBot         *telebot.Bot
	StatManager         *StatManager
	UserAccountManager  *UserAccountManager
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewTelegramBotManager creates new TelegramBotManager.
// Using for configuring with wire.
func NewTelegramBotManager(
	config *config.Config,
	logger *utils.Logger,
	statManager *StatManager,
	userAccountManager *UserAccountManager,
	telegramBot *utils.TelegramBot,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*TelegramBotManager, func(), error) {
	tbm := &TelegramBotManager{
		Config:              config,
		Logger:              logger,
		logger:              logger.Logger.Named("TelegramBotManager"),
		TelegramBot:         telegramBot,
		StatManager:         statManager,
		UserAccountManager:  userAccountManager,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("TelegramBotManager"),
	}

	ew := tools.GetErrorWrapper("NewTelegramBotManager")

	err := tbm.createBot()
	if err != nil {
		return nil, nil, ew(err)
	}

	go tbm.telegramBot.Start()

	return tbm, func() { tbm.telegramBot.Stop() }, nil
}

func (t *TelegramBotManager) createBot() error {
	ew := t.ErrorWrapperCreator.GetMethodWrapper("createBot")

	bot, err := t.TelegramBot.CreateBotDefault()
	if err != nil {
		return ew(err)
	}

	t.telegramBot = bot

	return nil
}

// GetBot returns [utils.TelegramBot] instance.
// If bot is not created, it will be created, but will panic if error occurred.
func (t *TelegramBotManager) GetBot() *telebot.Bot {
	if t.telegramBot == nil {
		tools.PanicOnError(t.createBot())
	}

	return t.telegramBot
}
