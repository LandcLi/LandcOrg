package v1

import (
	"github.com/LandcLi/landc-go/frame/pkg/meta"
)

// --- Organization Role Binding API ---

// AssignDepartmentRolesRequest binds roles to a department.
type AssignDepartmentRolesRequest struct {
	meta.Meta    `path:"/api/org-role/assign-department-roles" method:"POST"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	RoleIDs      []uint `json:"role_ids" binding:"required"`
}

type AssignDepartmentRolesResponse struct {
	Success bool `json:"success"`
}

// RevokeDepartmentRolesRequest unbinds roles from a department.
type RevokeDepartmentRolesRequest struct {
	meta.Meta    `path:"/api/org-role/revoke-department-roles" method:"POST"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	RoleIDs      []uint `json:"role_ids" binding:"required"`
}

type RevokeDepartmentRolesResponse struct {
	Success bool `json:"success"`
}

// GetDepartmentRolesRequest gets roles bound to a department.
type GetDepartmentRolesRequest struct {
	meta.Meta    `path:"/api/org-role/get-department-roles" method:"POST"`
	DepartmentID uint `json:"department_id" binding:"required"`
}

type GetDepartmentRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}

// AssignPositionRolesRequest binds roles to a position.
type AssignPositionRolesRequest struct {
	meta.Meta  `path:"/api/org-role/assign-position-roles" method:"POST"`
	PositionID uint   `json:"position_id" binding:"required"`
	RoleIDs    []uint `json:"role_ids" binding:"required"`
}

type AssignPositionRolesResponse struct {
	Success bool `json:"success"`
}

// RevokePositionRolesRequest unbinds roles from a position.
type RevokePositionRolesRequest struct {
	meta.Meta  `path:"/api/org-role/revoke-position-roles" method:"POST"`
	PositionID uint   `json:"position_id" binding:"required"`
	RoleIDs    []uint `json:"role_ids" binding:"required"`
}

type RevokePositionRolesResponse struct {
	Success bool `json:"success"`
}

// GetPositionRolesRequest gets roles bound to a position.
type GetPositionRolesRequest struct {
	meta.Meta  `path:"/api/org-role/get-position-roles" method:"POST"`
	PositionID uint `json:"position_id" binding:"required"`
}

type GetPositionRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}

// GetUserInheritedRolesRequest gets all roles a user inherits from org structure.
type GetUserInheritedRolesRequest struct {
	meta.Meta `path:"/api/org-role/get-user-inherited-roles" method:"POST"`
	UserID    uint `json:"user_id" binding:"required"`
}

type GetUserInheritedRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}
