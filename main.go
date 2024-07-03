// Package main - contains usage example of go-initial-project
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/telebot.v3"

	"github.com/roman-kart/go-initial-project/v2/project"
	"github.com/roman-kart/go-initial-project/v2/project/config"
	"github.com/roman-kart/go-initial-project/v2/project/managers"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

func main() {
	test()
}

func configureApp(app *project.Application) error {
	app.Logger.Logger.Info("Starting application")

	helpAdditionalMessage := "Чтобы получить справку по конкретной команде: `/help <команда без слэша>`"
	app.TelegramBotManager.AddCommonCommandsHandlers(&managers.CommonBotCommandsConfig{
		Start: managers.StartCommandConfig{
			Enabled: true,
			Message: "Привет!",
		},
		Help: managers.HelpCommandConfig{
			Enabled: true,
			MainHelpMessage: "Тестовый бот\n" +
				"\n" +
				helpAdditionalMessage,
			CommandsHelpMessages: map[string]managers.HelpCommandMessages{
				"/help": {
					ShortMessage: "Получить справку по боту",
					DetailMessage: "Напишите /help, чтобы посмотреть все команды.\n" +
						"\n" +
						helpAdditionalMessage,
				},
				"/start": {
					ShortMessage:  "Начать работу с ботом",
					DetailMessage: "Напишите /start",
				},
			},
		},
	})

	err := configureBot(app)

	return err
}

func configureBot(app *project.Application) error {
	botManager := app.TelegramBotManager
	bot := botManager.GetBot()
	adminsOnlyGroup := bot.Group()
	adminsOnlyGroup.Use(botManager.GetAdminOnlyMiddleware())
	adminsOnlyGroup.Handle("/admins_s3", func(c telebot.Context) error {
		ew := botManager.ErrorWrapperCreator.GetMethodWrapper("/admins_s3")

		if len(c.Args()) == 0 {
			return ew(c.Send("Command is not specified"))
		}

		command := c.Args()[0]
		switch command {
		case "list":
			objs, err := app.S3Manager.ListObjects(managers.ListObjectsInput{})
			if err != nil {
				return ew(c.Send("Error while listing objects"))
			}

			for _, obj := range objs {
				err = ew(
					c.Send(
						fmt.Sprintf(
							"Key: %s \nClass: %s \nSize: %d \nLastModifiued: %s",
							*obj.Key,
							obj.StorageClass,
							*obj.Size,
							obj.LastModified.Format("2006-01-02 15:04:05"),
						),
					),
				)
				if err != nil {
					return err
				}
			}

			return nil
		default:
			return ew(c.Send("Unknown command"))
		}
	})

	return nil
}

func test() {
	config.CountdownSecondsCount = 1

	rootPath := tools.GetRootPath()
	configFolder := rootPath + string(os.PathSeparator) + "project" + string(os.PathSeparator) + "config"
	app, cleanup, err := project.InitializeApplication(configFolder)

	defer cleanup()

	tools.PanicOnError(err)

	err = configureApp(app)
	tools.PanicOnError(err)

	readInput(app.Logger.Logger)
}

func readInput(logger *zap.Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			logger.Info("Program finished")
			break
		}
		// Help message
		logger.Info("For finish program enter `exit`")
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error while reading input", zap.Error(err))
	}
}
