package use_cases

import (
	"hole/entities"
)

type ItemRepository interface {
	Create(item *entities.Item) error
	FindByOwnerID(ownerID uint) ([]*entities.Item, error)
	// FindByIDAndOwner(id, ownerID uint) (*entities.Item, error)
	Update(id uint, name, desc string) error
	Delete(id uint) error
	ListItem() ([]*entities.Item, error)
}

type ItemUseCase struct {
	repo ItemRepository
}

func NewItemUseCase(repo ItemRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo}
}

func (uc *ItemUseCase) CreateItem(name, desc string) error {
	return uc.repo.Create(&entities.Item{
		ProductName: name,
		ProductDesc: desc,
	})
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
