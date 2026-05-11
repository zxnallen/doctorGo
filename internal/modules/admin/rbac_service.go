package admin

import "errors"

var ErrRoleNotFound = errors.New("role not found")

type RBACService struct {
	repo *Repository
}

func NewRBACService(repo *Repository) *RBACService {
	return &RBACService{repo: repo}
}

func (s *RBACService) ListRoles() ([]Role, error) {
	return s.repo.ListRoles()
}

func (s *RBACService) ListPermissions() ([]Permission, error) {
	return s.repo.ListPermissions()
}

func (s *RBACService) CreateRole(req CreateRoleRequest) (*Role, error) {
	role := &Role{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
	}
	if err := s.repo.CreateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RBACService) UpdateRole(id uint64, req UpdateRoleRequest) (*Role, error) {
	role, err := s.repo.FindRoleByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	role.Name = req.Name
	role.Status = req.Status
	if err := s.repo.UpdateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RBACService) SetRolePermissions(roleID uint64, permissionIDs []uint64) error {
	role, err := s.repo.FindRoleByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}
	return s.repo.SetRolePermissions(roleID, permissionIDs)
}

func (s *RBACService) SetAdminRoles(adminID uint64, roleIDs []uint64) error {
	return s.repo.SetAdminRoles(adminID, roleIDs)
}
