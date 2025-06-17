package services

import (
	"Netlfy/models"
	"Netlfy/repositories"
)

type ProfileService interface {
	FollowUser(currentUserID uint, username string) (*models.User, bool, error)
	UnfollowUser(currentUserID uint, username string) (*models.User, bool, error)
	GetProfile(username string, currentUserID *uint) (*models.User, bool, error)
}

type ProfileServiceImpl struct {
	userRepo   repositories.UserRepository
	followRepo repositories.FollowRepository
}

func NewProfileService(userRepo repositories.UserRepository, followRepo repositories.FollowRepository) ProfileService {
	return &ProfileServiceImpl{userRepo, followRepo}
}

func (s *ProfileServiceImpl) GetProfile(username string, currentUserID *uint) (*models.User, bool, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, false, err
	}
	var following bool
	if currentUserID != nil {
		following, _ = s.followRepo.IsFollowing(*currentUserID, user.ID)
	}
	return user, following, nil
}

func (s *ProfileServiceImpl) FollowUser(currentUserID uint, username string) (*models.User, bool, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, false, err
	}
	if err := s.followRepo.Follow(currentUserID, user.ID); err != nil {
		return nil, false, err
	}
	return user, true, nil
}

func (s *ProfileServiceImpl) UnfollowUser(currentUserID uint, username string) (*models.User, bool, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, false, err
	}
	if err := s.followRepo.Unfollow(currentUserID, user.ID); err != nil {
		return nil, false, err
	}
	return user, false, nil
}
