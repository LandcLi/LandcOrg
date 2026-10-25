package org_role

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/org_role/v1"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// OrgRoleController defines the organization-role binding management interface.
type OrgRoleController interface {
	AssignDepartmentRoles(ctx context.Context, req *v1.AssignDepartmentRolesRequest) (*v1.AssignDepartmentRolesResponse, error)
	RevokeDepartmentRoles(ctx context.Context, req *v1.RevokeDepartmentRolesRequest) (*v1.RevokeDepartmentRolesResponse, error)
	GetDepartmentRoles(ctx context.Context, req *v1.GetDepartmentRolesRequest) (*v1.GetDepartmentRolesResponse, error)
	AssignPositionRoles(ctx context.Context, req *v1.AssignPositionRolesRequest) (*v1.AssignPositionRolesResponse, error)
	RevokePositionRoles(ctx context.Context, req *v1.RevokePositionRolesRequest) (*v1.RevokePositionRolesResponse, error)
	GetPositionRoles(ctx context.Context, req *v1.GetPositionRolesRequest) (*v1.GetPositionRolesResponse, error)
	GetUserInheritedRoles(ctx context.Context, req *v1.GetUserInheritedRolesRequest) (*v1.GetUserInheritedRolesResponse, error)
}

var OrgRoleGateway = di.NewGateway[OrgRoleController]("org_role.controller")

func GetOrgRoleController() OrgRoleController {
	return OrgRoleGateway.Get()
}
