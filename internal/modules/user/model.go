package user

import "time"

type User struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Phone        string     `gorm:"size:32;index" json:"phone"`
	Email        string     `gorm:"size:128;index" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	AvatarURL    string     `gorm:"size:512" json:"avatar_url"`
	Status       int        `gorm:"not null;default:1" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
