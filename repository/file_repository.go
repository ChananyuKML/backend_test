package repository

import (
	"context"
	"hole/entities"
	"io"

	"github.com/minio/minio-go/v7"
)

type minioRepo struct {
	client     *minio.Client
	bucketName string
}

type FileRepository interface {
	Upload(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (minio.UploadInfo, error)
	GetObject(ctx context.Context, fileName string) (*entities.FileStream, error)
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

func (r *minioRepo) GetObject(ctx context.Context, fileName string) (*entities.FileStream, error) {
	object, err := r.client.GetObject(ctx, r.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		object.Close() // ALWAYS close if stat fails to free the connection
		return nil, err
	}

	return &entities.FileStream{
		Reader:      object, // Fiber's SendStream will close this
		ContentType: stat.ContentType,
		Size:        stat.Size,
	}, nil
}
