package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader struct {
	client *s3.Client
	bucket string
	// baseURL is used to construct the public URL after upload.
	// e.g. "https://my-bucket.s3.eu-central-1.amazonaws.com"
	baseURL string
}

func NewUploader(client *s3.Client, bucket, baseURL string) *Uploader {
	return &Uploader{client: client, bucket: bucket, baseURL: baseURL}
}

func (u *Uploader) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("s3.Upload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", u.baseURL, key)
	return url, nil
}
