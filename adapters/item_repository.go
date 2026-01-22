package adapters

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

func (r *ItemRepositoryPostgres) FindByIDAndOwner(id, ownerID uint) (*entities.Item, error) {
	var item entities.Item
	err := r.db.
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepositoryPostgres) FindByOwnerID(ownerID uint) ([]*entities.Item, error) {
	var items []*entities.Item
	err := r.db.Where("owner_id = ?", ownerID).Find(&items).Error
	return items, err
}

func (r *ItemRepositoryPostgres) Update(item *entities.Item) error {
	return r.db.
		Where("id = ? AND owner_id = ?", item.ID, item.OwnerID).
		Updates(map[string]interface{}{
			"description": item.Description,
		}).Error
}

func (r *ItemRepositoryPostgres) Delete(id, ownerID uint) error {
	result := r.db.
		Where("id = ? AND owner_id = ?", id, ownerID).
		Delete(&entities.Item{})

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return result.Error
}
