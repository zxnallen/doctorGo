package upload

type SignedURLRequest struct {
	BizType  string `json:"biz_type" binding:"required,max=64"`
	FileName string `json:"file_name" binding:"required,max=255"`
	MimeType string `json:"mime_type" binding:"required,max=128"`
	Size     int64  `json:"size" binding:"required,min=1"`
}

type SignedURLResponse struct {
	FileKey   string `json:"file_key"`
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
}
