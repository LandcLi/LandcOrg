package org_role

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterOrgRoleService(&orgRoleServiceImpl{}) }

type orgRoleServiceImpl struct{}

func (s *orgRoleServiceImpl) AssignDepartmentRoles(departmentID uint, roleIDs []uint) error {
	return dao.GetDepartmentRoleDAO().Assign(departmentID, roleIDs)
}

func (s *orgRoleServiceImpl) RevokeDepartmentRoles(departmentID uint, roleIDs []uint) error {
	return dao.GetDepartmentRoleDAO().Revoke(departmentID, roleIDs)
}

func (s *orgRoleServiceImpl) GetDepartmentRoles(departmentID uint) ([]uint, error) {
	return dao.GetDepartmentRoleDAO().GetRoleIDsByDepartmentID(departmentID)
}

func (s *orgRoleServiceImpl) AssignPositionRoles(positionID uint, roleIDs []uint) error {
	return dao.GetPositionRoleDAO().Assign(positionID, roleIDs)
}

func (s *orgRoleServiceImpl) RevokePositionRoles(positionID uint, roleIDs []uint) error {
	return dao.GetPositionRoleDAO().Revoke(positionID, roleIDs)
}

func (s *orgRoleServiceImpl) GetPositionRoles(positionID uint) ([]uint, error) {
	return dao.GetPositionRoleDAO().GetRoleIDsByPositionID(positionID)
}

func (s *orgRoleServiceImpl) GetUserInheritedRoleIDs(userID uint) ([]uint, error) {
	// Get all memberships of the user
	members, err := dao.GetMemberDAO().GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []uint{}, nil
	}

	// Collect department IDs and position IDs
	deptIDs := make([]uint, 0, len(members))
	posIDs := make([]uint, 0, len(members))
	for _, m := range members {
		deptIDs = append(deptIDs, m.DepartmentID)
		if m.PositionID > 0 {
			posIDs = append(posIDs, m.PositionID)
		}
	}

	// Get roles from departments
	deptRoleIDs, err := dao.GetDepartmentRoleDAO().GetRoleIDsByDepartmentIDs(deptIDs)
	if err != nil {
		return nil, err
	}

	// Get roles from positions
	posRoleIDs, err := dao.GetPositionRoleDAO().GetRoleIDsByPositionIDs(posIDs)
	if err != nil {
		return nil, err
	}

	// Merge and deduplicate
	roleSet := make(map[uint]struct{})
	for _, id := range deptRoleIDs {
		roleSet[id] = struct{}{}
	}
	for _, id := range posRoleIDs {
		roleSet[id] = struct{}{}
	}

	result := make([]uint, 0, len(roleSet))
	for id := range roleSet {
		result = append(result, id)
	}
	return result, nil
}
