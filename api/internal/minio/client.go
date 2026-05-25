package minio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	minioSDK "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	clientInstance *minioSDK.Client
	initError      error
	once           sync.Once
)

func NewClient() (*minioSDK.Client, error) {

	once.Do(func() {
		endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
		bucketName := strings.TrimSpace(os.Getenv("MINIO_BUCKET_NAME"))
		accessKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_USER"))
		secretKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD"))

		if endpoint == "" || bucketName == "" || accessKey == "" || secretKey == "" {
			initError = errors.New("MINIO_ENDPOINT, MINIO_BUCKET_NAME, MINIO_ROOT_USER, and MINIO_ROOT_PASSWORD must be set")
			return
		}

		useSSL := strings.TrimSpace(strings.ToLower(os.Getenv("MINIO_USE_SSL"))) == "true"

		client, err := minioSDK.New(endpoint, &minioSDK.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			initError = err
			return
		}

		ctx := context.Background()
		exists, err := client.BucketExists(ctx, bucketName)
		if err != nil {
			initError = fmt.Errorf("failed to verify bucket %q: %w", bucketName, err)
			return
		}
		if !exists {
			initError = fmt.Errorf("bucket %q does not exist", bucketName)
			return
		}

		clientInstance = client
	})

	if initError != nil {
		return nil, initError
	}
	return clientInstance, nil
}
