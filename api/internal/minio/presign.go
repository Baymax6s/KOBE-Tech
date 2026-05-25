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

func GeneratePresignedGetURL(
	client *minioSDK.Client,
	objectName string,
) (string, error) {

	u, err := client.PresignedGetObject(
		context.Background(),
		BucketName,
		objectName,
		10*time.Minute,
		url.Values{},
	)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}