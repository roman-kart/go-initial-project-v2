package managers

import (
	"errors"
	"fmt"
	"strings"

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

// GetDefaultSendOptions returns default send options.
func (t *TelegramBotManager) GetDefaultSendOptions() *telebot.SendOptions {
	return &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	}
}

// StartCommandConfig contains configurations for start command.
type StartCommandConfig struct {
	Enabled bool
	Message string
}

// HelpCommandMessages contains configurations for help command.
type HelpCommandMessages struct {
	ShortMessage  string
	DetailMessage string
}

// HelpCommandConfig contains configuration of one command's help message.
type HelpCommandConfig struct {
	Enabled         bool
	MainHelpMessage string
	// CommandsHelpMessages key is a handler endpoint with leading slash.
	CommandsHelpMessages map[string]HelpCommandMessages
}

// CommonBotCommandsConfig contains configurations for common bot commands.
type CommonBotCommandsConfig struct {
	Start StartCommandConfig
	Help  HelpCommandConfig
}

// AddCommonCommandsHandlers adds handlers for common bot commands.
func (t *TelegramBotManager) AddCommonCommandsHandlers(cfg *CommonBotCommandsConfig) {
	if cfg.Start.Enabled {
		ew := t.ErrorWrapperCreator.GetMethodWrapper("start_handler")

		t.GetBot().Handle("/start", func(c telebot.Context) error {
			message, err := TelegramStartCommandResponse(&cfg.Start)
			if err != nil {
				return ew(err)
			}

			return ew(c.Send(message, t.GetDefaultSendOptions()))
		})
	}

	if cfg.Help.Enabled {
		ew := t.ErrorWrapperCreator.GetMethodWrapper("help_handler")

		t.GetBot().Handle("/help", func(c telebot.Context) error {
			args := c.Args()

			message, err := TelegramHelpCommandResponse(&cfg.Help, args)
			if err != nil {
				return ew(err)
			}

			return ew(c.Send(message, t.GetDefaultSendOptions()))
		})
	}
}

// ErrNoMessage error if no message provided for command.
var ErrNoMessage = errors.New("no message")

// TelegramStartCommandResponse greeting user when user press Start button.
func TelegramStartCommandResponse(cfg *StartCommandConfig) (string, error) {
	ew := tools.GetErrorWrapper("TelegramStartCommandResponse")

	if cfg.Message == "" {
		return "", ew(ErrNoMessage)
	}

	return cfg.Message, nil
}

// TelegramHelpCommandResponse returns help message for all commands or concrete command.
// Commands will be sorted by command name.
func TelegramHelpCommandResponse(cfg *HelpCommandConfig, args []string) (string, error) {
	ew := tools.GetErrorWrapper("TelegramHelpCommandResponse")

	if cfg.MainHelpMessage == "" {
		return "", ew(ErrNoMessage)
	}

	if len(args) == 0 {
		commandsListMessagePart := ""

		commandKeysSorted := tools.SortMapKeys(cfg.CommandsHelpMessages)
		for _, commandName := range commandKeysSorted {
			commandConfig := cfg.CommandsHelpMessages[commandName]
			commandsListMessagePart += fmt.Sprintf("%s - %s\n", commandName, commandConfig.ShortMessage)
		}

		commandsListMessagePart = strings.TrimRight(commandsListMessagePart, "\n") // remove last \n

		if commandsListMessagePart == "" {
			return cfg.MainHelpMessage, nil
		}

		finalMessage := fmt.Sprintf("%s\n\n*Команды:*\n%s", cfg.MainHelpMessage, commandsListMessagePart)

		return finalMessage, nil
	}

	commandName := args[0]
	commandName = strings.TrimSpace(commandName)

	// remove leading slash if it exists for make sure only one slash will exists
	commandName = "/" + strings.TrimLeft(commandName, "/")

	commandConfig, ok := cfg.CommandsHelpMessages[commandName]

	if !ok {
		return fmt.Sprintf("Команда `%s` не найдена", commandName), nil
	}

	return fmt.Sprintf("%s\n\n%s", commandConfig.ShortMessage, commandConfig.DetailMessage), nil
}
