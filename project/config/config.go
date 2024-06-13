package config

import (
	"os"

	"github.com/jinzhu/configor"
)

type Config struct {
	ConfigFolder string `yaml:"-"`
	Clickhouse   struct {
		Host               string `yaml:"host" default:"localhost"`
		Port               int    `yaml:"port" default:"9000"`
		User               string `yaml:"user" default:"default"`
		Password           string `yaml:"password" default:""`
		Database           string `yaml:"database" default:"default"`
		IsNeedToRecreate   bool   `yaml:"is_need_to_recreate" default:"false"`
		AutoMigrate        bool   `yaml:"auto_migrate" default:"false"`
		IsNeedToInitialize bool   `yaml:"is_need_to_initialize" default:"false"`
		ConnMaxLifetime    int64  `yaml:"conn_max_lifetime" default:"60"`  // seconds
		ConnMaxIdleTime    int64  `yaml:"conn_max_idle_time" default:"60"` // seconds
		MaxIdleConns       int    `yaml:"max_idle_conns" default:"10"`
		MaxOpenConns       int    `yaml:"max_open_conns" default:"10"`
	} `yaml:"clickhouse"`
	Logger struct {
		Level string `yaml:"level" default:"info"`
	}
	Postgresql struct {
		Host               string `yaml:"host" default:"localhost"`
		Port               int    `yaml:"port" default:"5432"`
		User               string `yaml:"user" default:"postgres"`
		Password           string `yaml:"password" default:"postgres"`
		Database           string `yaml:"database" default:"lucky-gamer"`
		IsNeedToRecreate   bool   `yaml:"is_need_to_recreate" default:"false"`
		AutoMigrate        bool   `yaml:"auto_migrate" default:"false"`
		IsNeedToInitialize bool   `yaml:"is_need_to_initialize" default:"false"`
		ConnMaxLifetime    int64  `yaml:"conn_max_lifetime" default:"60"`  // seconds
		ConnMaxIdleTime    int64  `yaml:"conn_max_idle_time" default:"60"` // seconds
		MaxIdleConns       int    `yaml:"max_idle_conns" default:"10"`
		MaxOpenConns       int    `yaml:"max_open_conns" default:"10"`
	} `yaml:"postgresql"`
	IsDebug  bool `yaml:"is_debug" default:"false"`
	Telegram struct {
		Token  string  `yaml:"token"`
		ChatId int64   `yaml:"chat_id"`
		Admins []int64 `yaml:"admins"`
	} `yaml:"telegram"`
	RabbitMQ struct {
		Host     string `yaml:"host" default:"localhost"`
		Port     int    `yaml:"port" default:"5672"`
		User     string `yaml:"user" default:"guest"`
		Password string `yaml:"password" default:"guest"`
	} `yaml:"rabbitmq"`
	S3 struct {
		// paths from root
		ConfigPaths      []string `yaml:"config_paths"`
		CredentialsPaths []string `yaml:"credentials_paths"`
	} `yaml:"s3"`
}

func NewConfig(configFolder string) (*Config, error) {
	configPath := configFolder +
		string(os.PathSeparator) + "main.yaml"
	localConfigPath := configFolder +
		string(os.PathSeparator) + "main-local.yaml"
	var configPaths []string
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// only main config
		configPaths = append(configPaths, configPath)
	} else {
		// main and local
		configPaths = append(configPaths, localConfigPath, configPath)
	}
	config := Config{
		ConfigFolder: configFolder,
	}
	err := configor.Load(&config, configPaths...)
	return &config, err
}
