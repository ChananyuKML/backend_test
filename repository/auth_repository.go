package repository

import (
	"fmt"
	"hole/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryPostgres struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db}
}

func (r *UserRepositoryPostgres) Create(user *entities.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	fmt.Printf("hased password :%s", string(hash))

	return r.db.Create(&entities.User{
		Email:    user.Email,
		Password: string(hash),
	}).Error
}

func (r *UserRepositoryPostgres) FindByEmail(email string) (*entities.User, error) {
	var u entities.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &entities.User{ID: u.ID, Email: u.Email, Password: u.Password}, nil
}

func (r *UserRepositoryPostgres) FindByID(id uint) (*entities.User, error) {
	var u entities.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &entities.User{ID: u.ID, Email: u.Email}, nil
}
