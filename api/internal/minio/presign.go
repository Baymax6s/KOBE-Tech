package minio

import (
	"context"
	"net/url"
	"time"

	minioSDK "github.com/minio/minio-go/v7"
)

func GeneratePresignedPutURL(
	ctx context.Context,
	client *minioSDK.Client,
	bucketName string,
	objectName string,
) (string, error) {
	url, err := client.PresignedPutObject(
		ctx,
		bucketName,
		objectName,
		10*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func GeneratePresignedGetURL(
	ctx context.Context,
	client *minioSDK.Client,
	bucketName string,
	objectName string,
) (string, error) {

	u, err := client.PresignedGetObject(
		ctx,
		bucketName,
		objectName,
		10*time.Minute,
		url.Values{},
	)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
