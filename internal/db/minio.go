package db

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioBlobStore struct {
	client *minio.Client
}

const (
	linkExpiresTime = time.Minute * 5
)

func NewMinioBlobStore(endpoint, accessKey, secretKey string, useSSL bool) (*MinioBlobStore, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioBlobStore{client: client}, nil
}

func (m *MinioBlobStore) UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader, size int64, contentType string) error {
	ok, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !ok {
		err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	_, err = m.client.PutObject(ctx, bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (m *MinioBlobStore) GetReader(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (m *MinioBlobStore) GetAccessLink(ctx context.Context, bucketName, objectName string) (*url.URL, error) {
	presignedURL, err := m.client.PresignedGetObject(context.Background(), bucketName, objectName, linkExpiresTime, url.Values{})
	if err != nil {
		return &url.URL{}, err
	}

	return presignedURL, nil
}

func (m *MinioBlobStore) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
