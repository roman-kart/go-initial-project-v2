package managers

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// S3Manager is a struct for managing files in systems like S3.
type S3Manager struct {
	Config              *config.Config
	Logger              *utils.Logger
	logger              *zap.Logger
	ErrorWrapperCreator tools.ErrorWrapperCreator
	S3Client            *utils.S3
}

// NewS3Manager creates a new instance of S3Manager.
// Using for configuring with wire.
func NewS3Manager(
	config *config.Config,
	logger *utils.Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
	s3Client *utils.S3,
) (*S3Manager, error) {
	s3Manager := &S3Manager{
		Config:              config,
		Logger:              logger,
		logger:              logger.Logger.Named("S3Manager"),
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("S3Manager"),
		S3Client:            s3Client,
	}

	ew := tools.GetErrorWrapper("NewS3Manager")

	client, err := s3Manager.GetClient()
	if err != nil {
		return nil, ew(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(s3Manager.Config.S3Manager.Timeout)*time.Second,
	)
	defer cancel()

	// Testing connection to S3.
	_, err = client.ListBuckets(ctx, nil)
	if err != nil {
		return nil, ew(err)
	}

	return s3Manager, nil
}

// GetClient creates a new S3 client.
func (s3Manager *S3Manager) GetClient() (*s3.Client, error) {
	ew := s3Manager.ErrorWrapperCreator.GetMethodWrapper("GetClient")

	client, err := s3Manager.S3Client.GetClient()
	if err != nil {
		return nil, ew(err)
	}

	return client, nil
}

// ListObjectsInput argument for [S3Manager.ListObjects] function.
type ListObjectsInput struct {
	Bucket  string
	MaxKeys int32
	Prefix  string
}

// ListObjects lists all objects in a default bucket.
func (s3Manager *S3Manager) ListObjects(input ListObjectsInput) ([]types.Object, error) {
	ew := s3Manager.ErrorWrapperCreator.GetMethodWrapper("ListObjects")

	client, err := s3Manager.GetClient()
	if err != nil {
		return nil, ew(err)
	}

	objectsList := []types.Object{}

	var continuationToken *string

	bucket := tools.FirstNonEmpty(input.Bucket, s3Manager.Config.S3Manager.Bucket)
	maxKeys := tools.FirstNonEmpty(input.MaxKeys, s3Manager.Config.S3Manager.MaxKeys)

	for {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(s3Manager.Config.S3Manager.Timeout)*time.Second,
		)
		defer cancel()

		result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            &bucket,
			MaxKeys:           &maxKeys,
			ContinuationToken: continuationToken,
			Prefix:            &input.Prefix,
		})
		if err != nil {
			return objectsList, ew(err)
		}

		objectsList = append(objectsList, result.Contents...)

		if result.IsTruncated != nil && *result.IsTruncated {
			continuationToken = result.NextContinuationToken
		} else {
			break
		}
	}

	return objectsList, nil
}
