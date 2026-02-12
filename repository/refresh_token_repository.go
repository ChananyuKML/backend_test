package repository

import (
	"hole/entities"

	"gorm.io/gorm"
)

type RefreshTokenRepositoryPostgres struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepositoryPostgres {
	return &RefreshTokenRepositoryPostgres{db}
}

func (r *RefreshTokenRepositoryPostgres) Create(t *entities.RefreshToken) error {
	return r.db.Create(&entities.RefreshToken{
		UserID:    t.UserID,
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
	}).Error
}

func (r *RefreshTokenRepositoryPostgres) FindByToken(token string) (*entities.RefreshToken, error) {
	var t entities.RefreshToken
	if err := r.db.Where("token = ?", token).First(&t).Error; err != nil {
		return nil, err
	}
	return &entities.RefreshToken{
		ID:        t.ID,
		UserID:    t.UserID,
		Token:     t.Token,
		Revoked:   t.Revoked,
		ExpiresAt: t.ExpiresAt,
	}, nil
}

func (r *RefreshTokenRepositoryPostgres) Revoke(token string) error {
	return r.db.Model(&entities.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}
