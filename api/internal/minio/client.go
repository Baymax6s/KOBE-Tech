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
	bucketName     string
	initError      error
	once           sync.Once
)

// BucketName は NewClient が解決・検証済みのバケット名を返す。
// 環境変数を読み直すと正規化（前後空白の除去）がズレる余地があるため、
// バケット名の出どころはこの 1 か所に集約する。NewClient 成功後にのみ有効。
func BucketName() string {
	return bucketName
}

func NewClient() (*minioSDK.Client, error) {

	once.Do(func() {
		endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
		bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET_NAME"))
		accessKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_USER"))
		secretKey := strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD"))

		if endpoint == "" || bucket == "" || accessKey == "" || secretKey == "" {
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
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			initError = fmt.Errorf("failed to verify bucket %q: %w", bucket, err)
			return
		}
		if !exists {
			initError = fmt.Errorf("bucket %q does not exist", bucket)
			return
		}

		clientInstance = client
		bucketName = bucket
	})

	if initError != nil {
		return nil, initError
	}
	return clientInstance, nil
}
