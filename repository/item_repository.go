package repository

import (
	"errors"
	"hole/entities"

	"gorm.io/gorm"
)

type ItemRepositoryPostgres struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepositoryPostgres {
	return &ItemRepositoryPostgres{db}
}

func (r *ItemRepositoryPostgres) Create(item *entities.Item) error {
	return r.db.Create(item).Error
}

func (r *ItemRepositoryPostgres) FindByIDAndOwner(id uint) (*entities.Item, error) {
	var item entities.Item
	err := r.db.
		Where("product_id = ?", id).
		First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepositoryPostgres) ListItem() ([]*entities.Item, error) {
	var items []*entities.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *ItemRepositoryPostgres) FindByOwnerID(ownerID uint) ([]*entities.Item, error) {
	var items []*entities.Item
	err := r.db.Where("owner_id = ?", ownerID).Find(&items).Error
	return items, err
}

func (r *ItemRepositoryPostgres) Update(id uint, name string, desc string) error {

	result := r.db.Model(&entities.Item{}).
		Where("product_id = ?", id).
		Updates(entities.Item{
			ProductName: name,
			ProductDesc: desc,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (r *ItemRepositoryPostgres) Delete(id uint) error {
	result := r.db.Where("product_id = ?", id).Delete(&entities.Item{})

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}
