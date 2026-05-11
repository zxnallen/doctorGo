package upload

import "time"

type File struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	BizType   string    `gorm:"size:64;index" json:"biz_type"`
	FileKey   string    `gorm:"size:512;not null;uniqueIndex" json:"file_key"`
	URL       string    `gorm:"size:1024;not null" json:"url"`
	MimeType  string    `gorm:"size:128" json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedBy uint64    `gorm:"index" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (File) TableName() string {
	return "upload_files"
}
