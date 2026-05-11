package news

import (
	"errors"

	"gorm.io/gorm"

	"doctor-go/internal/infrastructure/mysql"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *mysql.DB) *Repository {
	return &Repository{db: db.DB}
}

func (r *Repository) List(req ListRequest, offset int, limit int) ([]News, int64, error) {
	var items []News
	var total int64

	query := r.db.Model(&News{}).Where("status = ?", 1)
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("published_at DESC, id DESC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) AdminList(req AdminListRequest, offset int, limit int) ([]News, int64, error) {
	var items []News
	var total int64

	query := r.db.Model(&News{})
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) FindPublishedByID(id uint64) (*News, error) {
	var item News
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindByID(id uint64) (*News, error) {
	var item News
	err := r.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) Create(item *News) error {
	return r.db.Create(item).Error
}

func (r *Repository) Update(item *News) error {
	return r.db.Save(item).Error
}

func (r *Repository) Delete(id uint64) error {
	return r.db.Delete(&News{}, id).Error
}

func (r *Repository) IncrementViewCount(id uint64) error {
	return r.db.Model(&News{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *Repository) ListCategories(onlyEnabled bool) ([]Category, error) {
	var items []Category
	query := r.db.Model(&Category{})
	if onlyEnabled {
		query = query.Where("status = ?", 1)
	}
	err := query.Order("sort DESC, id DESC").Find(&items).Error
	return items, err
}

func (r *Repository) FindCategoryByID(id uint64) (*Category, error) {
	var item Category
	err := r.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) CreateCategory(item *Category) error {
	return r.db.Create(item).Error
}

func (r *Repository) UpdateCategory(item *Category) error {
	return r.db.Save(item).Error
}

func (r *Repository) DeleteCategory(id uint64) error {
	return r.db.Delete(&Category{}, id).Error
}
