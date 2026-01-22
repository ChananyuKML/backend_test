package use_cases

import (
	"hole/entities"
)

type ItemRepository interface {
	Create(item *entities.Item) error
	FindByOwnerID(ownerID uint) ([]*entities.Item, error)
	FindByIDAndOwner(id, ownerID uint) (*entities.Item, error)
	Update(item *entities.Item) error
	Delete(id, ownerID uint) error
}

type ItemUseCase struct {
	repo ItemRepository
}

func NewItemUseCase(repo ItemRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo}
}

func (uc *ItemUseCase) CreateItem(ownerID uint, desc string) error {
	return uc.repo.Create(&entities.Item{
		OwnerID:     ownerID,
		Description: desc,
	})
}

func (uc *ItemUseCase) GetMyItems(ownerID uint) ([]*entities.Item, error) {
	return uc.repo.FindByOwnerID(ownerID)
}

func (uc *ItemUseCase) UpdateItem(id, ownerID uint, desc string) error {
	return uc.repo.Update(&entities.Item{
		ID:          id,
		OwnerID:     ownerID,
		Description: desc,
	})
}

func (uc *ItemUseCase) DeleteItem(id, ownerID uint) error {
	return uc.repo.Delete(id, ownerID)
}
