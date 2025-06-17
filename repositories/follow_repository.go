package repositories

import (
	"Netlfy/database"
	"Netlfy/models"
	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(followerID, followeeID uint) error
	Unfollow(followerID, followeeID uint) error
	IsFollowing(followerID, followeeID uint) (bool, error)
}

type FollowRepositoryImpl struct{}

func NewFollowRepository() FollowRepository {
	return &FollowRepositoryImpl{}
}

func (r *FollowRepositoryImpl) Follow(followerID, followeeID uint) error {
	follow := models.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	return database.DB.FirstOrCreate(&follow, follow).Error
}

func (r *FollowRepositoryImpl) Unfollow(followerID, followeeID uint) error {
	return database.DB.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{}).Error
}

func (r *FollowRepositoryImpl) IsFollowing(followerID, followeeID uint) (bool, error) {
	var follow models.Follow
	err := database.DB.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).First(&follow).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
