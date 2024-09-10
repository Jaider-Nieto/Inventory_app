package repository

import (
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindAllUsers() ([]models.User, error) {
	var users []models.User

	err := r.DB.Find(&users).Error

	return users, err
}
func (r *UserRepository) FindUserByID(id string) (models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error

	return user, err
}
func (r *UserRepository) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error

	return user, err
}
func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	err := r.DB.Create(&user).Error

	return user, err
}
func (r *UserRepository) DeleteUser(id string) error {
	err := r.DB.Delete(&models.User{}, id).Error
	return err
}
func (r *UserRepository) UpdateUser(user models.User) error {
	err := r.DB.Save(&user).Error
	return err
}
