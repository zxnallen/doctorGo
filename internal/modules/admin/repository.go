package admin

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"doctor-go/internal/infrastructure/mysql"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *mysql.DB) *Repository {
	return &Repository{db: db.DB}
}

func (r *Repository) FindByUsername(username string) (*AdminUser, error) {
	var admin AdminUser
	err := r.db.Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *Repository) Create(admin *AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *Repository) UpdateLastLogin(adminID uint64, at time.Time) error {
	return r.db.Model(&AdminUser{}).Where("id = ?", adminID).Update("last_login_at", at).Error
}

func (r *Repository) HasPermission(adminID uint64, permissionCode string) (bool, error) {
	var count int64
	err := r.db.Table("admin_permissions AS p").
		Joins("JOIN admin_role_permissions AS rp ON rp.permission_id = p.id").
		Joins("JOIN admin_roles AS r ON r.id = rp.role_id").
		Joins("JOIN admin_user_roles AS ur ON ur.role_id = r.id").
		Where("ur.admin_id = ? AND p.code = ? AND r.status = ?", adminID, permissionCode, 1).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ListRoles() ([]Role, error) {
	var items []Role
	err := r.db.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *Repository) ListPermissions() ([]Permission, error) {
	var items []Permission
	err := r.db.Order("resource ASC, action ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *Repository) CreateRole(role *Role) error {
	return r.db.Create(role).Error
}

func (r *Repository) UpdateRole(role *Role) error {
	return r.db.Save(role).Error
}

func (r *Repository) FindRoleByID(id uint64) (*Role, error) {
	var role Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) SetRolePermissions(roleID uint64, permissionIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			relation := RolePermission{RoleID: roleID, PermissionID: permissionID}
			if err := tx.Create(&relation).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) SetAdminRoles(adminID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("admin_id = ?", adminID).Delete(&AdminUserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			relation := AdminUserRole{AdminID: adminID, RoleID: roleID}
			if err := tx.Create(&relation).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) EnsureRBAC(adminID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		role := Role{Code: "super_admin", Name: "超级管理员", Status: 1}
		if err := tx.Where("code = ?", role.Code).FirstOrCreate(&role).Error; err != nil {
			return err
		}

		permissions := defaultPermissions()
		for _, permission := range permissions {
			p := permission
			if err := tx.Where("code = ?", p.Code).FirstOrCreate(&p).Error; err != nil {
				return err
			}
			relation := RolePermission{RoleID: role.ID, PermissionID: p.ID}
			if err := tx.Where("role_id = ? AND permission_id = ?", role.ID, p.ID).FirstOrCreate(&relation).Error; err != nil {
				return err
			}
		}

		if adminID > 0 {
			relation := AdminUserRole{AdminID: adminID, RoleID: role.ID}
			if err := tx.Where("admin_id = ? AND role_id = ?", adminID, role.ID).FirstOrCreate(&relation).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func defaultPermissions() []Permission {
	return []Permission{
		{Code: "news:list", Name: "资讯列表", Resource: "news", Action: "list"},
		{Code: "news:create", Name: "创建资讯", Resource: "news", Action: "create"},
		{Code: "news:update", Name: "更新资讯", Resource: "news", Action: "update"},
		{Code: "news:status", Name: "更新资讯状态", Resource: "news", Action: "change_status"},
		{Code: "news:delete", Name: "删除资讯", Resource: "news", Action: "delete"},
		{Code: "news_category:list", Name: "资讯分类列表", Resource: "news_category", Action: "list"},
		{Code: "news_category:create", Name: "创建资讯分类", Resource: "news_category", Action: "create"},
		{Code: "news_category:update", Name: "更新资讯分类", Resource: "news_category", Action: "update"},
		{Code: "news_category:status", Name: "更新资讯分类状态", Resource: "news_category", Action: "change_status"},
		{Code: "news_category:delete", Name: "删除资讯分类", Resource: "news_category", Action: "delete"},
		{Code: "operation_log:list", Name: "操作日志列表", Resource: "operation_log", Action: "list"},
		{Code: "rbac:list", Name: "查看角色权限", Resource: "rbac", Action: "list"},
		{Code: "rbac:manage", Name: "管理角色权限", Resource: "rbac", Action: "manage"},
	}
}
