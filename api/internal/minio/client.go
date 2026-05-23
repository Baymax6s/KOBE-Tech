package minio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	minioSDK "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	DefaultMinioEndpoint = "localhost:9000"
	BucketName            = "kobe-tech-avatar"
)

func NewClient() (*minioSDK.Client, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	if endpoint == "" {
		endpoint = DefaultMinioEndpoint
	}

	accessKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_USER"))
	secretKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD"))
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("MINIO_ROOT_USER and MINIO_ROOT_PASSWORD must be set")
	}

	client, err := minioSDK.New(endpoint, &minioSDK.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to verify bucket %q: %w", BucketName, err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket %q does not exist", BucketName)
	}

	return client, nil
}