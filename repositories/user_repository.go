package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(userID uint) (*models.User, error)
	UpdateUser(user *models.User) error
	GetByUsername(username string) (*models.User, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) GetByID(userID uint) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepositoryImpl) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
