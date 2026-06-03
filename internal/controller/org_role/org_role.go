package org_role

import (
	"context"

	orgRoleApi "github.com/LandcLi/LandcOrg/api/org_role"
	v1 "github.com/LandcLi/LandcOrg/api/org_role/v1"
	"github.com/LandcLi/LandcOrg/service"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

func init() { di.Provide[orgRoleApi.OrgRoleController]("org_role.controller", &orgRoleCtrl{}) }

type orgRoleCtrl struct{}

func (c *orgRoleCtrl) AssignDepartmentRoles(ctx context.Context, req *v1.AssignDepartmentRolesRequest) (*v1.AssignDepartmentRolesResponse, error) {
	if err := service.GetOrgRoleService().AssignDepartmentRoles(req.DepartmentID, req.RoleIDs); err != nil {
		return nil, err
	}
	return &v1.AssignDepartmentRolesResponse{Success: true}, nil
}

func (c *orgRoleCtrl) RevokeDepartmentRoles(ctx context.Context, req *v1.RevokeDepartmentRolesRequest) (*v1.RevokeDepartmentRolesResponse, error) {
	if err := service.GetOrgRoleService().RevokeDepartmentRoles(req.DepartmentID, req.RoleIDs); err != nil {
		return nil, err
	}
	return &v1.RevokeDepartmentRolesResponse{Success: true}, nil
}

func (c *orgRoleCtrl) GetDepartmentRoles(ctx context.Context, req *v1.GetDepartmentRolesRequest) (*v1.GetDepartmentRolesResponse, error) {
	roleIDs, err := service.GetOrgRoleService().GetDepartmentRoles(req.DepartmentID)
	if err != nil {
		return nil, err
	}
	return &v1.GetDepartmentRolesResponse{RoleIDs: roleIDs}, nil
}

func (c *orgRoleCtrl) AssignPositionRoles(ctx context.Context, req *v1.AssignPositionRolesRequest) (*v1.AssignPositionRolesResponse, error) {
	if err := service.GetOrgRoleService().AssignPositionRoles(req.PositionID, req.RoleIDs); err != nil {
		return nil, err
	}
	return &v1.AssignPositionRolesResponse{Success: true}, nil
}

func (c *orgRoleCtrl) RevokePositionRoles(ctx context.Context, req *v1.RevokePositionRolesRequest) (*v1.RevokePositionRolesResponse, error) {
	if err := service.GetOrgRoleService().RevokePositionRoles(req.PositionID, req.RoleIDs); err != nil {
		return nil, err
	}
	return &v1.RevokePositionRolesResponse{Success: true}, nil
}

func (c *orgRoleCtrl) GetPositionRoles(ctx context.Context, req *v1.GetPositionRolesRequest) (*v1.GetPositionRolesResponse, error) {
	roleIDs, err := service.GetOrgRoleService().GetPositionRoles(req.PositionID)
	if err != nil {
		return nil, err
	}
	return &v1.GetPositionRolesResponse{RoleIDs: roleIDs}, nil
}

func (c *orgRoleCtrl) GetUserInheritedRoles(ctx context.Context, req *v1.GetUserInheritedRolesRequest) (*v1.GetUserInheritedRolesResponse, error) {
	roleIDs, err := service.GetOrgRoleService().GetUserInheritedRoleIDs(req.UserID)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInheritedRolesResponse{RoleIDs: roleIDs}, nil
}
