package upload

import (
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"doctor-go/internal/infrastructure/oss"
)

var ErrOSSNotConfigured = errors.New("oss not configured")

type Service struct {
	oss  *oss.Client
	repo *Repository
}

func NewService(ossClient *oss.Client, repo *Repository) *Service {
	return &Service{oss: ossClient, repo: repo}
}

func (s *Service) CreateSignedURL(userID uint64, req SignedURLRequest) (*SignedURLResponse, error) {
	if s.oss == nil || !s.oss.IsConfigured() {
		return nil, ErrOSSNotConfigured
	}

	ext := strings.ToLower(filepath.Ext(req.FileName))
	fileKey := req.BizType + "/" + time.Now().Format("20060102") + "/" + uuid.NewString() + ext
	uploadURL, err := s.oss.SignedPutURL(fileKey, 900)
	if err != nil {
		return nil, err
	}
	publicURL := s.oss.PublicURL(fileKey)

	if err := s.repo.Create(&File{
		BizType:   req.BizType,
		FileKey:   fileKey,
		URL:       publicURL,
		MimeType:  req.MimeType,
		Size:      req.Size,
		CreatedBy: userID,
	}); err != nil {
		return nil, err
	}

	return &SignedURLResponse{
		FileKey:   fileKey,
		UploadURL: uploadURL,
		PublicURL: publicURL,
	}, nil
}
