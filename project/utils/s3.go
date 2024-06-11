package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	c "github.com/roman-kart/go-initial-project/project/config"
)

type S3 struct {
	Config     *c.Config
	Logger     *Logger
	Postgresql *Postgresql
}

func NewS3(Config *c.Config, Logger *Logger, Postgresql *Postgresql) *S3 {
	return &S3{
		Config:     Config,
		Logger:     Logger,
		Postgresql: Postgresql,
	}
}

func (s *S3) GetClient() (*s3.Client, error) {
	var configPaths []string
	for _, path := range s.Config.S3.ConfigPaths {
		configPaths = append(configPaths, GetPathFromRoot(path))
	}

	var credentialsPaths []string
	for _, path := range s.Config.S3.CredentialsPaths {
		credentialsPaths = append(credentialsPaths, GetPathFromRoot(path))
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigFiles(configPaths),
		config.WithSharedCredentialsFiles(credentialsPaths),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
