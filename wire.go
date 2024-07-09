//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/roman-kart/go-initial-project/v2/project/managers"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
	"github.com/roman-kart/go-initial-project/v2/project/utils"
)

func InitializeApplication(
	configFolder string,
	configCountdownSecondsCount uint,
) (*Application, func(), error) {
	wire.Build(
		NewConfig,

		NewS3ManagerConfig,
		NewTelegramBotManagerConfig,
		NewClickHouseConfig,
		NewLoggerConfig,
		NewPostgresqlConfig,
		NewRabbitMQConfig,
		NewS3Config,
		NewTelegramConfig,

		tools.NewErrorWrapperCreator,
		utils.NewClickHouse,
		utils.NewLogger,
		utils.NewPostgresql,
		utils.NewRabbitMQ,
		utils.NewS3,
		utils.NewTelegram,

		managers.NewStatManager,
		managers.NewTelegramBotManager,
		managers.NewUserAccountManager,
		managers.NewS3Manager,

		NewApplication,
	)
	return &Application{}, func() {}, nil
}
