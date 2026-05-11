package upload

import (
	"gorm.io/gorm"

	"doctor-go/internal/infrastructure/mysql"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *mysql.DB) *Repository {
	return &Repository{db: db.DB}
}

func (r *Repository) Create(file *File) error {
	return r.db.Create(file).Error
}
