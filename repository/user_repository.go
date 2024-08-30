package repository

import (
	"github.com/jaider-nieto/ecommerce-go/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindUsers() ([]models.User, error) {
	var err error
	return []models.User{}, err
}

func (r *UserRepository) FindUserById(id int) (models.User, error) {
	var err error
	return models.User{}, err
}

func (r *UserRepository) CreateUserById(models.User) (models.User, error) {
	var err error
	return models.User{}, err
}

func (r *UserRepository) DeleteUserById(id int) error {
	var err error
	return err
}

func (r *UserRepository) UpdateUserById(models.User) (models.User, error) {
	var err error
	return models.User{}, err
}
