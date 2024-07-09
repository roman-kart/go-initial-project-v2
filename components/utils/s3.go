package utils

import (
	"context"
	"go.uber.org/zap"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/roman-kart/go-initial-project/v2/components/tools"
)

type S3Config struct {
	ConfigPaths      []string
	CredentialsPaths []string
	ConfigFolder     string
}

// S3 manipulates connections to services like Amazon S3.
type S3 struct {
	Config              *S3Config
	logger              *zap.Logger
	Postgresql          *Postgresql
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewS3 creates an instance of S3.
// Using for configuring with wire.
func NewS3(
	config *S3Config,
	logger *zap.Logger,
	postgresql *Postgresql,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *S3 {
	return &S3{
		Config:              config,
		logger:              logger,
		Postgresql:          postgresql,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("S3"),
	}
}

// GetClient creates a client without caching.
func (s *S3) GetClient() (*s3.Client, error) {
	ew := s.ErrorWrapperCreator.GetMethodWrapper("GetClient")

	configPaths := []string{}
	for _, path := range s.Config.ConfigPaths {
		configPaths = append(configPaths, s.Config.ConfigFolder+string(os.PathSeparator)+path)
	}

	credentialsPaths := []string{}
	for _, path := range s.Config.CredentialsPaths {
		credentialsPaths = append(credentialsPaths, s.Config.ConfigFolder+string(os.PathSeparator)+path)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigFiles(configPaths),
		config.WithSharedCredentialsFiles(credentialsPaths),
	)
	if err != nil {
		return nil, ew(err)
	}

	return s3.NewFromConfig(cfg), nil
}
