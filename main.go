// Package main - contains usage example of go-initial-project
package main

import (
	"os"
	"time"

	"github.com/roman-kart/go-initial-project/project"
	"github.com/roman-kart/go-initial-project/project/managers"
	"github.com/roman-kart/go-initial-project/project/tools"
)

func main() {
	test()
}

func test() {
	rootPath := tools.GetRootPath()
	configFolder := rootPath + string(os.PathSeparator) + "project" + string(os.PathSeparator) + "config"
	app, cleanup, err := project.InitializeApplication(configFolder)

	defer cleanup()

	tools.PanicOnError(err)

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
				"/start": {
					ShortMessage:  "Начать работу с ботом",
					DetailMessage: "Напишите /start",
				},
				"/help": {
					ShortMessage: "Получить справку по боту",
					DetailMessage: "Напишите /help, чтобы посмотреть все команды.\n" +
						"\n" +
						helpAdditionalMessage,
				},
			},
		},
	})

	time.Sleep(1 * time.Minute)
}
