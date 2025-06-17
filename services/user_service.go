package services

import (
	"Netlfy/models"
	"Netlfy/repositories"
	"Netlfy/utils"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, email, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	GetUserByID(userID uint) (*models.User, error)
	UpdateUser(userID uint, email, username, password, image, bio *string) (*models.User, string, error)
}

type UserServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImpl{repo}
}

func (s *UserServiceImpl) Register(username, email, password string) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *UserServiceImpl) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("Wrong password or Email")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *UserServiceImpl) GetUserByID(userID uint) (*models.User, error) {
	return s.repo.GetByID(userID)
}

func (s *UserServiceImpl) UpdateUser(userID uint, email, username, password, image, bio *string) (*models.User, string, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, "", err
	}

	if email != nil {
		user.Email = *email
	}
	if username != nil {
		user.Username = *username
	}
	if password != nil {
		hashedPassword, err := utils.HashPassword(*password)
		if err != nil {
			return nil, "", err
		}
		user.Password = hashedPassword
	}
	if image != nil {
		user.Image = *image
	}
	if bio != nil {
		user.Bio = *bio
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
