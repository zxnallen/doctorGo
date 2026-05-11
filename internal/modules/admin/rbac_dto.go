package admin

type CreateRoleRequest struct {
	Code   string `json:"code" binding:"required,max=64"`
	Name   string `json:"name" binding:"required,max=64"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name" binding:"required,max=64"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

type SetRolePermissionsRequest struct {
	PermissionIDs []uint64 `json:"permission_ids" binding:"required"`
}

type SetAdminRolesRequest struct {
	RoleIDs []uint64 `json:"role_ids" binding:"required"`
}
