package admin

import "time"

type AdminUser struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Nickname     string     `gorm:"size:64;not null;default:''" json:"nickname"`
	AvatarURL    string     `gorm:"size:512;not null;default:''" json:"avatar_url"`
	Status       int        `gorm:"not null;default:1" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}

type Role struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "admin_roles"
}

type Permission struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Resource  string    `gorm:"size:64;not null;default:''" json:"resource"`
	Action    string    `gorm:"size:64;not null;default:''" json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "admin_permissions"
}

type AdminUserRole struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	AdminID   uint64    `gorm:"uniqueIndex:uk_admin_role;not null" json:"admin_id"`
	RoleID    uint64    `gorm:"uniqueIndex:uk_admin_role;not null" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (AdminUserRole) TableName() string {
	return "admin_user_roles"
}

type RolePermission struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	RoleID       uint64    `gorm:"uniqueIndex:uk_role_permission;not null" json:"role_id"`
	PermissionID uint64    `gorm:"uniqueIndex:uk_role_permission;not null" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "admin_role_permissions"
}
