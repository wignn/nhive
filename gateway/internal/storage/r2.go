package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/novelhive/gateway/internal/config"
)

type R2Client struct {
	S3Client *s3.Client
	Config   *cfg.Config
}

func NewR2Client(c *cfg.Config) (*R2Client, error) {
	if (c.R2AccountID == "" && c.R2Endpoint == "") || c.R2AccessKeyID == "" || c.R2SecretAccessKey == "" || c.R2BucketName == "" {
		return nil, fmt.Errorf("R2 credentials not fully configured")
	}

	endpointURL := c.R2Endpoint
	if endpointURL == "" {
		endpointURL = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.R2AccountID)
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpointURL,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.R2AccessKeyID, c.R2SecretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return &R2Client{
		S3Client: client,
		Config:   c,
	}, nil
}

func (r *R2Client) UploadImage(ctx context.Context, file multipart.File, filename, contentType string) (string, string, error) {
	_, err := r.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.Config.R2BucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload image: %v", err)
	}

	baseURL := r.Config.R2PublicURL
	if baseURL == "" {
		// Fallback: This URL requires AWS SigV4 auth and won't work in a browser directly.
		// The user MUST configure Public Access in Cloudflare and set R2_PUBLIC_URL.
		baseURL = fmt.Sprintf("%s/%s", r.Config.R2Endpoint, r.Config.R2BucketName)
	}

	return filename, baseURL, nil
}
