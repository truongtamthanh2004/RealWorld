package models

type Follow struct {
	ID         uint `gorm:"primaryKey"`
	FollowerID uint
	FolloweeID uint
}
