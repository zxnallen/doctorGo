package admin

import (
	"time"

	"doctor-go/internal/pkg/pagination"

	"gorm.io/gorm"

	"doctor-go/internal/infrastructure/mysql"
)

type OperationLog struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	AdminID    uint64         `gorm:"index;not null" json:"admin_id"`
	Username   string         `gorm:"size:64;not null;default:''" json:"username"`
	Action     string         `gorm:"size:64;not null;index" json:"action"`
	Resource   string         `gorm:"size:64;not null;index" json:"resource"`
	ResourceID uint64         `gorm:"index" json:"resource_id"`
	Method     string         `gorm:"size:16;not null;default:''" json:"method"`
	Path       string         `gorm:"size:255;not null;default:''" json:"path"`
	StatusCode int            `gorm:"not null;default:0" json:"status_code"`
	Success    bool           `gorm:"not null;default:true" json:"success"`
	Remark     string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OperationLog) TableName() string {
	return "admin_operation_logs"
}

type OperationLogRepository struct {
	db *gorm.DB
}

type OperationLogListRequest struct {
	Page     int
	Size     int
	AdminID  *uint64
	Action   string
	Resource string
}

type OperationLogListResponse struct {
	Items []OperationLog `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

func NewOperationLogRepository(db *mysql.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db.DB}
}

func (r *OperationLogRepository) Create(log *OperationLog) error {
	return r.db.Create(log).Error
}

func (r *OperationLogRepository) List(req OperationLogListRequest) ([]OperationLog, int64, error) {
	var items []OperationLog
	var total int64
	query := r.db.Model(&OperationLog{})
	if req.AdminID != nil {
		query = query.Where("admin_id = ?", *req.AdminID)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Resource != "" {
		query = query.Where("resource = ?", req.Resource)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := pagination.Page{Page: req.Page, Size: req.Size}.Normalize()
	err := query.Order("id DESC").Offset(page.Offset()).Limit(page.Size).Find(&items).Error
	return items, total, err
}
