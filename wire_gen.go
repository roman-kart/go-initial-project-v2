// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/roman-kart/go-initial-project/v2/components/managers"
	"github.com/roman-kart/go-initial-project/v2/components/tools"
	"github.com/roman-kart/go-initial-project/v2/components/utils"
)

// Injectors from wire.go:

func InitializeApplication(configFolder string, configCountdownSecondsCount uint) (*Application, func(), error) {
	config, err := NewConfig(configFolder, configCountdownSecondsCount)
	if err != nil {
		return nil, nil, err
	}
	clickHouseConfig := NewClickHouseConfig(config)
	loggerConfig := NewLoggerConfig(config)
	logger, cleanup, err := utils.NewLogger(loggerConfig)
	if err != nil {
		return nil, nil, err
	}
	errorWrapperCreator := tools.NewErrorWrapperCreator()
	clickHouse, cleanup2, err := utils.NewClickHouse(clickHouseConfig, logger, errorWrapperCreator)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	postgresqlConfig := NewPostgresqlConfig(config)
	postgresql, cleanup3, err := utils.NewPostgresql(postgresqlConfig, logger, errorWrapperCreator)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	rabbitMQConfig := NewRabbitMQConfig(config)
	rabbitMQ := utils.NewRabbitMQ(rabbitMQConfig, logger, errorWrapperCreator)
	s3Config := NewS3Config(config)
	s3 := utils.NewS3(s3Config, logger, postgresql, errorWrapperCreator)
	telegramConfig := NewTelegramConfig(config)
	telegramBot := utils.NewTelegram(telegramConfig, logger, errorWrapperCreator)
	statManager, err := managers.NewStatManager(logger, clickHouse, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	telegramBotManagerConfig := NewTelegramBotManagerConfig(config)
	userAccountManager, err := managers.NewUserAccountManager(logger, postgresql, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	telegramBotManager, cleanup4, err := managers.NewTelegramBotManager(telegramBotManagerConfig, logger, statManager, userAccountManager, telegramBot, errorWrapperCreator)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	s3ManagerConfig := NewS3ManagerConfig(config)
	s3Manager, err := managers.NewS3Manager(s3ManagerConfig, logger, errorWrapperCreator, s3)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(config, clickHouse, logger, postgresql, rabbitMQ, s3, telegramBot, statManager, telegramBotManager, userAccountManager, s3Manager)
	return application, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}