package service

import "github.com/LandcLi/landc-go/frame/pkg/di"

// OrgRoleService defines the organization-driven role binding interface.
type OrgRoleService interface {
	// AssignDepartmentRoles binds roles to a department.
	AssignDepartmentRoles(departmentID uint, roleIDs []uint) error
	// RevokeDepartmentRoles unbinds roles from a department.
	RevokeDepartmentRoles(departmentID uint, roleIDs []uint) error
	// GetDepartmentRoles gets all role IDs bound to a department.
	GetDepartmentRoles(departmentID uint) ([]uint, error)

	// AssignPositionRoles binds roles to a position.
	AssignPositionRoles(positionID uint, roleIDs []uint) error
	// RevokePositionRoles unbinds roles from a position.
	RevokePositionRoles(positionID uint, roleIDs []uint) error
	// GetPositionRoles gets all role IDs bound to a position.
	GetPositionRoles(positionID uint) ([]uint, error)

	// GetUserInheritedRoleIDs gets all role IDs a user inherits from departments and positions.
	GetUserInheritedRoleIDs(userID uint) ([]uint, error)
}

func RegisterOrgRoleService(impl OrgRoleService) {
	di.Provide[OrgRoleService]("org_role.service", impl)
}

func GetOrgRoleService() OrgRoleService {
	return di.Require[OrgRoleService]("org_role.service")
}
