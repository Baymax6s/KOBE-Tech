package minio

import (
	"context"
	"net/url"
	"time"

	minioSDK "github.com/minio/minio-go/v7"
)

func GeneratePresignedPutURL(
	client *minioSDK.Client,
	objectName string,
) (string, error) {
	url, err := client.PresignedPutObject(
		context.Background(),
		BucketName,
		objectName,
		10*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

// 💡 ✨ここを追加：閲覧（GET）用の署名付きURLを発行する関数
func GeneratePresignedGetURL(
	client *minioSDK.Client,
	objectName string,
) (string, error) {
	// MinIO SDKの PresignedGetObject を呼び出します（有効期限はとりあえず同じく10分）
	u, err := client.PresignedGetObject(
		context.Background(),
		BucketName,
		objectName,
		10*time.Minute,
		url.Values{}, // GET時の追加クエリパラメータ（今回は空でOK）
	)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}