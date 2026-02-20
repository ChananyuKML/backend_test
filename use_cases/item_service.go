package use_cases

import (
	"context"
	"fmt"
	"hole/entities"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type ItemRepository interface {
	Create(item *entities.Item) error
	FindByOwnerID(ownerID uint) ([]*entities.Item, error)
	// FindByIDAndOwner(id, ownerID uint) (*entities.Item, error)
	Update(id uint, name, desc string) error
	Delete(id uint) error
	ListItem() ([]*entities.Item, error)
}

type FileRepository interface {
	Upload(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (minio.UploadInfo, error)
	GetObject(ctx context.Context, fileName string) (io.Reader, error)
}

type ItemUseCase struct {
	repo     ItemRepository
	fileRepo FileRepository
}

func NewItemUseCase(repo ItemRepository, fileRepo FileRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo, fileRepo: fileRepo}
}

func (uc *ItemUseCase) CreateItem(name, desc string, ctx context.Context, imageKey string) error {
	item := &entities.Item{
		ProductName:     name,
		ProductDesc:     desc,
		ProductImageKey: imageKey, // e.g., "products-images/177...jpg"
	}

	return uc.repo.Create(item)
}

func (uc *ItemUseCase) GetMyItems(ownerID uint) ([]*entities.Item, error) {
	return uc.repo.FindByOwnerID(ownerID)
}

func (uc *ItemUseCase) GetAllItems() ([]*entities.Item, error) {
	return uc.repo.ListItem()
}

func (uc *ItemUseCase) UpdateItem(id uint, name, desc string) error {
	return uc.repo.Update(id, name, desc)
}

func (uc *ItemUseCase) DeleteItem(id uint) error {
	return uc.repo.Delete(id)
}

func (u *ItemUseCase) UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (string, error) {
	fileName := fmt.Sprintf("products-images/%d.jpg", time.Now().UnixNano())
	info, err := u.fileRepo.Upload(ctx, fileName, file, size, contentType)

	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (u *ItemUseCase) GetImageStream(ctx context.Context, imageKey string) (io.Reader, error) {
	if imageKey == "" {
		return nil, fmt.Errorf("empty image key")
	}
	return u.fileRepo.GetObject(ctx, imageKey)
}
