package persistence

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client defines an interface for getting objects from S3.
type S3Client interface {
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
}

// s3Client is a concrete implementation of S3Client.
type s3Client struct {
	client *s3.Client
}

// NewS3Client returns a new S3Client.
func NewS3Client(cfg aws.Config) S3Client {
	return &s3Client{
		client: s3.NewFromConfig(cfg),
	}
}

// GetObject retrieves the S3 object for given bucket and key.
func (s *s3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
