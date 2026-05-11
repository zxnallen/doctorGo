package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/modules/user"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *mysql.DB) *Repository {
	return &Repository{db: db.DB}
}

func (r *Repository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) FindByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) UpdateLastLogin(userID uint64, at time.Time) error {
	return r.db.Model(&user.User{}).Where("id = ?", userID).Update("last_login_at", at).Error
}
