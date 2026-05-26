package minio

import (
	"context"
	"log"
)

// ResolveAvatarURL はアバターがアップロード済みのときだけ presigned GET URL を返す。
// ストレージ接続や署名に失敗してもプロフィール/記事の取得自体は止めず、空文字を返して
// フロント側のフォールバック表示（既定アイコン）に委ねる。
func ResolveAvatarURL(ctx context.Context, objectKey string, isUploaded bool) string {
	if !isUploaded || objectKey == "" {
		return ""
	}

	client, err := NewClient()
	if err != nil {
		log.Printf("avatar url: connect storage: %v", err)
		return ""
	}

	url, err := GeneratePresignedGetURL(ctx, client, bucketName, objectKey)
	if err != nil {
		log.Printf("avatar url: presign: %v", err)
		return ""
	}

	return url
}
