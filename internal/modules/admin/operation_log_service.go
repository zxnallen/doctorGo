package admin

import (
	"time"
)

type OperationLogService struct {
	repo *OperationLogRepository
}

func NewOperationLogService(repo *OperationLogRepository) *OperationLogService {
	return &OperationLogService{repo: repo}
}

func (s *OperationLogService) Create(log *OperationLog) error {
	return s.repo.Create(log)
}

func (s *OperationLogService) List(req OperationLogListRequest) ([]OperationLog, int64, error) {
	return s.repo.List(req)
}

func NewOperationLog(adminID uint64, username string, action string, resource string, resourceID uint64, method string, path string, statusCode int, success bool, remark string) *OperationLog {
	return &OperationLog{
		AdminID:    adminID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Success:    success,
		Remark:     remark,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
