package main

import (
	"context"
	"fmt"
	"github.com/roman-kart/go-initial-project/v2/components/managers"
	"github.com/roman-kart/go-initial-project/v2/components/tools"
	"github.com/roman-kart/go-initial-project/v2/components/utils"
	"os"
	"time"
)

// Config holds the application configuration.
type Config struct {
	ConfigFolder string `yaml:"-"`
	Clickhouse   struct {
		Host               string `default:"localhost" yaml:"host"`
		Port               int    `default:"9000"      yaml:"port"`
		User               string `default:"default"   yaml:"user"`
		Password           string `default:""          yaml:"password"`
		Database           string `default:"default"   yaml:"database"`
		IsNeedToRecreate   bool   `default:"false"     yaml:"is_need_to_recreate"`
		AutoMigrate        bool   `default:"false"     yaml:"auto_migrate"`
		IsNeedToInitialize bool   `default:"false"     yaml:"is_need_to_initialize"`
		ConnMaxLifetime    int64  `default:"60"        yaml:"conn_max_lifetime"`  // seconds
		ConnMaxIdleTime    int64  `default:"60"        yaml:"conn_max_idle_time"` // seconds
		MaxIdleConns       int    `default:"10"        yaml:"max_idle_conns"`
		MaxOpenConns       int    `default:"10"        yaml:"max_open_conns"`
	} `yaml:"clickhouse"`
	Logger struct {
		Level    string `default:"info" yaml:"level"`
		Sampling struct {
			Initial    int `default:"100" yaml:"initial"`
			Thereafter int `default:"200" yaml:"thereafter"`
		}
	}
	Postgresql struct {
		Host               string `default:"localhost"   yaml:"host"`
		Port               int    `default:"5432"        yaml:"port"`
		User               string `default:"postgres"    yaml:"user"`
		Password           string `default:"postgres"    yaml:"password"`
		Database           string `default:"lucky-gamer" yaml:"database"`
		IsNeedToRecreate   bool   `default:"false"       yaml:"is_need_to_recreate"`
		AutoMigrate        bool   `default:"false"       yaml:"auto_migrate"`
		IsNeedToInitialize bool   `default:"false"       yaml:"is_need_to_initialize"`
		ConnMaxLifetime    int64  `default:"60"          yaml:"conn_max_lifetime"`  // seconds
		ConnMaxIdleTime    int64  `default:"60"          yaml:"conn_max_idle_time"` // seconds
		MaxIdleConns       int    `default:"10"          yaml:"max_idle_conns"`
		MaxOpenConns       int    `default:"10"          yaml:"max_open_conns"`
	} `yaml:"postgresql"`
	IsDebug  bool `default:"false"   yaml:"is_debug"`
	Telegram struct {
		Token      string  `yaml:"token"`
		Admins     []int64 `yaml:"admins"`
		LongPoller struct {
			Timeout uint `default:"10" yaml:"timeout"`
		} `yaml:"long_poller"`
	} `yaml:"telegram"`
	RabbitMQ struct {
		Host     string `default:"localhost" yaml:"host"`
		Port     int    `default:"5672"      yaml:"port"`
		User     string `default:"guest"     yaml:"user"`
		Password string `default:"guest"     yaml:"password"`
	} `yaml:"rabbitmq"`
	S3 struct {
		// paths from root
		ConfigPaths      []string `yaml:"config_paths"`
		CredentialsPaths []string `yaml:"credentials_paths"`
	} `yaml:"s3"`
	S3Manager struct {
		Timeout uint   `default:"10"   yaml:"timeout"`
		Bucket  string `yaml:"bucket"`
		MaxKeys int32  `default:"1000" yaml:"max_keys"`
	} `yaml:"s3_manager"`
}

// NewConfig creates a new config.
// Using for read configuration from config files.
func NewConfig(configFolder string, configCountdownSecondsCount uint) (*Config, error) {
	configPath := configFolder +
		string(os.PathSeparator) + "main.yaml"
	localConfigPath := configFolder +
		string(os.PathSeparator) + "main-local.yaml"

	configPaths := []string{
		localConfigPath,
		configPath,
	}

	config := Config{
		ConfigFolder: configFolder,
	}

	err := tools.LoadConfig(configPaths, &config)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: %w", err)
	}

	if config.IsDebug {
		fmt.Printf("Config loaded successfully!\n"+
			"Config:%+v\n", config,
		)
	}

	alertsForProperties := map[string]bool{
		"Enable recreation of clickhouse - TABLES WILL BE DELETED THAT CREATED":   config.Clickhouse.IsNeedToRecreate,
		"Enable auto migrate of clickhouse - TABLE WILL BE ALTERED AUTOMATICALLY": config.Clickhouse.AutoMigrate,
		"Enable recreation of postgresql - TABLES WILL BE DELETED THAT CREATED":   config.Postgresql.IsNeedToRecreate,
		"Enable auto migrate of postgresql - TABLE WILL BE ALTERED AUTOMATICALLY": config.Postgresql.AutoMigrate,
	}

	for message, needToDisplay := range alertsForProperties {
		if needToDisplay {
			tools.RedOutputCmd(message)
		}
	}

	if configCountdownSecondsCount > 0 {
		tools.CountdownCmd(context.Background(), "CHECK CONFIG", time.Second, configCountdownSecondsCount)
	}

	return &config, err
}

func NewS3ManagerConfig(config *Config) *managers.S3ManagerConfig {
	return &managers.S3ManagerConfig{
		Timeout: config.S3Manager.Timeout,
		Bucket:  config.S3Manager.Bucket,
		MaxKeys: config.S3Manager.MaxKeys,
	}
}

func NewTelegramBotManagerConfig(config *Config) *managers.TelegramBotManagerConfig {
	return &managers.TelegramBotManagerConfig{
		Token:  config.Telegram.Token,
		Admins: config.Telegram.Admins,
		LongPoller: struct {
			Timeout uint
		}{
			Timeout: config.Telegram.LongPoller.Timeout,
		},
	}
}

func NewClickHouseConfig(config *Config) *utils.ClickHouseConfig {
	return &utils.ClickHouseConfig{
		Host:               config.Clickhouse.Host,
		Port:               config.Clickhouse.Port,
		User:               config.Clickhouse.User,
		Password:           config.Clickhouse.Password,
		Database:           config.Clickhouse.Database,
		IsNeedToRecreate:   config.Clickhouse.IsNeedToRecreate,
		AutoMigrate:        config.Clickhouse.AutoMigrate,
		IsNeedToInitialize: config.Clickhouse.IsNeedToInitialize,
		ConnMaxLifetime:    config.Clickhouse.ConnMaxLifetime,
		ConnMaxIdleTime:    config.Clickhouse.ConnMaxIdleTime,
		MaxIdleConns:       config.Clickhouse.MaxIdleConns,
		MaxOpenConns:       config.Clickhouse.MaxOpenConns,
		IsDebug:            config.IsDebug,
	}
}

func NewLoggerConfig(config *Config) *utils.LoggerConfig {
	return &utils.LoggerConfig{
		Level: config.Logger.Level,
		Sampling: struct {
			Initial    int
			Thereafter int
		}{
			Initial:    config.Logger.Sampling.Initial,
			Thereafter: config.Logger.Sampling.Thereafter,
		},
	}
}

func NewPostgresqlConfig(config *Config) *utils.PostgresqlConfig {
	return &utils.PostgresqlConfig{
		Host:               config.Postgresql.Host,
		Port:               config.Postgresql.Port,
		User:               config.Postgresql.User,
		Password:           config.Postgresql.Password,
		Database:           config.Postgresql.Database,
		IsNeedToRecreate:   config.Postgresql.IsNeedToRecreate,
		AutoMigrate:        config.Postgresql.AutoMigrate,
		IsNeedToInitialize: config.Postgresql.IsNeedToInitialize,
		ConnMaxLifetime:    config.Postgresql.ConnMaxLifetime,
		ConnMaxIdleTime:    config.Postgresql.ConnMaxIdleTime,
		MaxIdleConns:       config.Postgresql.MaxIdleConns,
		MaxOpenConns:       config.Postgresql.MaxOpenConns,
		IsDebug:            config.IsDebug,
	}
}

func NewRabbitMQConfig(config *Config) *utils.RabbitMQConfig {
	return &utils.RabbitMQConfig{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
	}
}

func NewS3Config(config *Config) *utils.S3Config {
	return &utils.S3Config{
		ConfigPaths:      config.S3.ConfigPaths,
		CredentialsPaths: config.S3.CredentialsPaths,
		ConfigFolder:     config.ConfigFolder,
	}
}

func NewTelegramConfig(config *Config) *utils.TelegramConfig {
	return &utils.TelegramConfig{
		LongPoller: struct {
			Timeout uint
		}{
			Timeout: config.Telegram.LongPoller.Timeout,
		},
	}
}
