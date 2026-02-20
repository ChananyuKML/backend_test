package repository

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type minioRepo struct {
	client     *minio.Client
	bucketName string
}

type FileRepository interface {
	Upload(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (minio.UploadInfo, error)
	GetObject(ctx context.Context, fileName string) (io.Reader, error)
}

func NewMinioRepo(client *minio.Client, bucket string) FileRepository {
	return &minioRepo{
		client:     client,
		bucketName: bucket,
	}
}

// minio_repo.go
func (r *minioRepo) Upload(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (minio.UploadInfo, error) {
	info, err := r.client.PutObject(ctx, r.bucketName, fileName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return info, err
}

func (r *minioRepo) GetObject(ctx context.Context, fileName string) (io.Reader, error) {
	// GetObject returns a stream of the file content
	object, err := r.client.GetObject(ctx, r.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
