package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/configor"
)

// CountdownSecondsCount is the number of seconds to countdown before config was loaded.
var CountdownSecondsCount uint //nolint:gochecknoglobals

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
// Using for configuring with wire.
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

	if err == nil {
		fmt.Printf("config loaded: %+v\n", config)

		alertsForProperties := map[string]bool{
			"Enable recreation of clickhouse - TABLES WILL BE DELETED THAT CREATED":   config.Clickhouse.IsNeedToRecreate,
			"Enable auto migrate of clickhouse - TABLE WILL BE ALTERED AUTOMATICALLY": config.Clickhouse.AutoMigrate,
			"Enable recreation of postgresql - TABLES WILL BE DELETED THAT CREATED":   config.Postgresql.IsNeedToRecreate,
			"Enable auto migrate of postgresql - TABLE WILL BE ALTERED AUTOMATICALLY": config.Postgresql.AutoMigrate,
		}

		for message, needToDisplay := range alertsForProperties {
			if needToDisplay {
				redOutput(message)
			}
		}

		countdownSecondsCount := uint(10) //nolint:mnd

		if CountdownSecondsCount > 0 {
			countdownSecondsCount = CountdownSecondsCount
		}

		countdown(context.Background(), "CHECK CONFIG", time.Second, countdownSecondsCount)
	} else {
		err = fmt.Errorf("NewConfig: %w", err)
	}

	return &config, err
}

func countdown(ctx context.Context, message string, delay time.Duration, count uint) {
	fmt.Println(message)
	fmt.Printf("Countdown: %d\n", count)

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for i := count; i > 0; i-- {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Printf("%d ", i)
		}
	}

	fmt.Println("0")
}

func redOutput(message string) {
	fmt.Println("\033[31m" + message + "\033[0m")
}
