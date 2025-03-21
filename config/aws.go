package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// NewAWSConfig loads and returns the default AWS config.
func NewAWSConfig(ctx context.Context) (aws.Config, error) {
	profile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_REGION")
	opts := []func(*config.LoadOptions) error{}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	} else {
		opts = append(opts, config.WithRegion("us-east-1"))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
