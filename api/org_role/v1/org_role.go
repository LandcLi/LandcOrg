package v1

// --- Organization Role Binding API ---

// AssignDepartmentRolesRequest binds roles to a department.
type AssignDepartmentRolesRequest struct {
	DepartmentID uint   `json:"department_id" binding:"required"`
	RoleIDs      []uint `json:"role_ids" binding:"required"`
}

type AssignDepartmentRolesResponse struct {
	Success bool `json:"success"`
}

// RevokeDepartmentRolesRequest unbinds roles from a department.
type RevokeDepartmentRolesRequest struct {
	DepartmentID uint   `json:"department_id" binding:"required"`
	RoleIDs      []uint `json:"role_ids" binding:"required"`
}

type RevokeDepartmentRolesResponse struct {
	Success bool `json:"success"`
}

// GetDepartmentRolesRequest gets roles bound to a department.
type GetDepartmentRolesRequest struct {
	DepartmentID uint `json:"department_id" binding:"required"`
}

type GetDepartmentRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}

// AssignPositionRolesRequest binds roles to a position.
type AssignPositionRolesRequest struct {
	PositionID uint   `json:"position_id" binding:"required"`
	RoleIDs    []uint `json:"role_ids" binding:"required"`
}

type AssignPositionRolesResponse struct {
	Success bool `json:"success"`
}

// RevokePositionRolesRequest unbinds roles from a position.
type RevokePositionRolesRequest struct {
	PositionID uint   `json:"position_id" binding:"required"`
	RoleIDs    []uint `json:"role_ids" binding:"required"`
}

type RevokePositionRolesResponse struct {
	Success bool `json:"success"`
}

// GetPositionRolesRequest gets roles bound to a position.
type GetPositionRolesRequest struct {
	PositionID uint `json:"position_id" binding:"required"`
}

type GetPositionRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}

// GetUserInheritedRolesRequest gets all roles a user inherits from org structure.
type GetUserInheritedRolesRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type GetUserInheritedRolesResponse struct {
	RoleIDs []uint `json:"role_ids"`
}
