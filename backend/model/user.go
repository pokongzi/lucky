package model

import "time"

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	OpenID    string `gorm:"uniqueIndex;size:64;not null"`
	Nickname  string `gorm:"size:64;not null"`
	AvatarURL string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
