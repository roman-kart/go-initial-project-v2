package environment

import (
	"errors"
	"github.com/roman-kart/go-initial-project/v2/project/config"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
	"github.com/roman-kart/go-initial-project/v2/project/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type InitializerConfig struct {
	RootPath string

	CreateAutocompleteShell   bool
	CreateGitignore           bool
	CreateGolangCIConfig      bool
	CreateHelperShell         bool
	CreateReadmeMd            bool
	CreateDefaultConfigFolder bool
}

type Initializer struct {
	Logger              *utils.Logger
	logger              *zap.Logger
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

var ErrRootPathIsNotADirectory = errors.New("root path is not a directory")

func NewInitializer(
	logger *utils.Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *Initializer {
	return &Initializer{
		Logger:              logger,
		logger:              logger.Logger.Named("Initializer"),
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("Initializer"),
	}
}

func (i *Initializer) Initialize(cfg InitializerConfig) error {
	ew := i.ErrorWrapperCreator.GetMethodWrapper("Initialize")
	logger := i.logger.Named("Initialize")

	if cfg.RootPath == "" {
		cfg.RootPath = tools.GetRootPath()
		logger.Info("Root path is empty, using default root path", zap.String("path", cfg.RootPath))
	}

	joinWithRoot := func(p ...string) string {
		return filepath.Join(cfg.RootPath, filepath.Join(p...))
	}

	f, err := os.Stat(cfg.RootPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error("Root path does not exist", zap.Error(err))
		}

		return ew(err)
	}

	if !f.IsDir() {
		logger.Error("Root path is not a directory", zap.Error(ErrRootPathIsNotADirectory))
		return ew(ErrRootPathIsNotADirectory)
	}

	logger.Info("Root path exists and is a directory", zap.String("path", cfg.RootPath))

	if cfg.CreateAutocompleteShell {
		err = os.WriteFile(joinWithRoot("autocomplete.sh"), []byte(AutocompleteShellScript), 0644)
		if err != nil {
			logger.Error("Failed to create autocomplete shell script", zap.Error(err))
			return ew(err)
		}
		logger.Info("Created autocomplete shell script")
	}

	if cfg.CreateGitignore {
		err = os.WriteFile(joinWithRoot(".gitignore"), []byte(Gitignore), 0644)
		if err != nil {
			logger.Error("Failed to create gitignore", zap.Error(err))
			return ew(err)
		}
		logger.Info("Created default .gitignore file")
	}

	if cfg.CreateGolangCIConfig {
		err = os.WriteFile(joinWithRoot(".golangci.yaml"), []byte(GolangCIConfig), 0644)
		if err != nil {
			logger.Error("Failed to create golangci-lint config", zap.Error(err))
			return ew(err)
		}
		logger.Info("Created default .golangci.yaml file")
	}

	if cfg.CreateHelperShell {
		err = os.WriteFile(joinWithRoot("helper.sh"), []byte(HelperShellScript), 0644)
		if err != nil {
			logger.Error("Failed to create helper shell", zap.Error(err))
			return ew(err)
		}
		logger.Info("Created helper shell")
	}

	if cfg.CreateReadmeMd {
		err = os.WriteFile(joinWithRoot("README.md"), []byte(ReadmeMd), 0644)
		if err != nil {
			logger.Error("Failed to create readme file", zap.Error(err))
			return ew(err)
		}
		logger.Info("Created default README.md file")
	}

	if cfg.CreateDefaultConfigFolder {
		configFolderPath := joinWithRoot("config")
		awsConfigPath := joinWithRoot("config/aws")

		// create folders

		err = os.Mkdir(configFolderPath, 0755)
		if err != nil {
			logger.Error("Failed to create default config folder", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default config folder")

		err = os.Mkdir(awsConfigPath, 0755)
		if err != nil {
			logger.Error("Failed to create default aws config folder", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default aws config folder")

		// fill aws folder

		err = os.WriteFile(filepath.Join(awsConfigPath, "config"), []byte(config.AwsConfigExample), 0644)
		if err != nil {
			logger.Error("Failed to create default aws config file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default aws config file")

		err = os.WriteFile(filepath.Join(awsConfigPath, "config.ex"), []byte(config.AwsConfigExample), 0644)
		if err != nil {
			logger.Error("Failed to create default aws config file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default aws config example file")

		err = os.WriteFile(filepath.Join(awsConfigPath, "credentials"), []byte(config.AwsCredentialsExample), 0644)
		if err != nil {
			logger.Error("Failed to create default aws credentials file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default aws credentials file")

		err = os.WriteFile(filepath.Join(awsConfigPath, "credentials.ex"), []byte(config.AwsCredentialsExample), 0644)
		if err != nil {
			logger.Error("Failed to create default aws credentials file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created default aws credentials example file")

		// fill config folder

		err = os.WriteFile(filepath.Join(configFolderPath, ".gitignore"), []byte(config.Gitignore), 0644)
		if err != nil {
			logger.Error("Failed to create default gitignore file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Created .gitignore file")

		err = os.WriteFile(filepath.Join(configFolderPath, "main.yaml"), []byte(config.MainConfig), 0644)
		if err != nil {
			logger.Error("Failed to create default main config file", zap.Error(err))
			return ew(err)
		}

		err = os.WriteFile(filepath.Join(configFolderPath, "main-local.yaml"), []byte(config.MainLocalConfigExample), 0644)
		if err != nil {
			logger.Error("Failed to create default main-local config file", zap.Error(err))
			return ew(err)
		}

		err = os.WriteFile(filepath.Join(configFolderPath, "main-local.yaml.ex"), []byte(config.MainLocalConfigExample), 0644)
		if err != nil {
			logger.Error("Failed to create default main-local example config file", zap.Error(err))
			return ew(err)
		}

		logger.Info("Config folder created and filled. Change data in files to yours")
	}

	return nil

}
