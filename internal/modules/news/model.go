package news

import "time"

type News struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null;index" json:"title"`
	Summary     string     `gorm:"size:500" json:"summary"`
	Content     string     `gorm:"type:longtext" json:"content"`
	CoverURL    string     `gorm:"size:512" json:"cover_url"`
	Author      string     `gorm:"size:64" json:"author"`
	CategoryID  uint64     `gorm:"index" json:"category_id"`
	Status      int        `gorm:"not null;default:1;index" json:"status"`
	ViewCount   int64      `gorm:"not null;default:0" json:"view_count"`
	PublishedAt *time.Time `gorm:"index" json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (News) TableName() string {
	return "news"
}

type Category struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Sort      int       `gorm:"not null;default:0" json:"sort"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Category) TableName() string {
	return "news_categories"
}
