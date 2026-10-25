package member

import (
	"fmt"
	"time"

	v1 "github.com/LandcLi/LandcOrg/api/member/v1"
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterMemberService(&memberServiceImpl{}) }

type memberServiceImpl struct{}

func (s *memberServiceImpl) AddMember(userID, departmentID, positionID uint, isLeader bool) error {
	// Check if already exists
	existing, _ := dao.GetMemberDAO().GetByUserAndDept(userID, departmentID)
	if existing != nil && existing.ID != 0 {
		return fmt.Errorf("user %d is already a member of department %d", userID, departmentID)
	}

	m := &model.OrgMember{
		UserID:       userID,
		DepartmentID: departmentID,
		PositionID:   positionID,
		IsLeader:     isLeader,
		JoinedAt:     time.Now(),
	}
	return dao.GetMemberDAO().Create(m)
}

func (s *memberServiceImpl) RemoveMember(userID, departmentID uint) error {
	return dao.GetMemberDAO().Delete(userID, departmentID)
}

func (s *memberServiceImpl) UpdateMember(userID, departmentID uint, positionID *uint, isLeader *bool) error {
	m, err := dao.GetMemberDAO().GetByUserAndDept(userID, departmentID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}
	if positionID != nil {
		m.PositionID = *positionID
	}
	if isLeader != nil {
		m.IsLeader = *isLeader
	}
	return dao.GetMemberDAO().Update(m)
}

func (s *memberServiceImpl) GetUserDepartments(userID uint) ([]v1.MemberInfo, error) {
	members, err := dao.GetMemberDAO().GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 批量预取部门与岗位，避免循环内逐个 GetByID 造成 N+1。
	deptIDs := make([]uint, 0, len(members))
	posIDs := make([]uint, 0, len(members))
	for _, m := range members {
		deptIDs = append(deptIDs, m.DepartmentID)
		if m.PositionID > 0 {
			posIDs = append(posIDs, m.PositionID)
		}
	}
	deptNames := departmentNameMap(deptIDs)
	posNames := positionNameMap(posIDs)

	result := make([]v1.MemberInfo, 0, len(members))
	for _, m := range members {
		info := v1.MemberInfo{
			UserID:       m.UserID,
			DepartmentID: m.DepartmentID,
			PositionID:   m.PositionID,
			IsLeader:     m.IsLeader,
			JoinedAt:     m.JoinedAt,
		}
		info.DepartmentName = deptNames[m.DepartmentID]
		if m.PositionID > 0 {
			info.PositionName = posNames[m.PositionID]
		}
		result = append(result, info)
	}
	return result, nil
}

func (s *memberServiceImpl) GetDepartmentMembers(departmentID uint, page, pageSize int) ([]v1.MemberInfo, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	members, total, err := dao.GetMemberDAO().GetByDepartmentID(departmentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 批量预取岗位名，消除 N+1。
	posIDs := make([]uint, 0, len(members))
	for _, m := range members {
		if m.PositionID > 0 {
			posIDs = append(posIDs, m.PositionID)
		}
	}
	posNames := positionNameMap(posIDs)

	result := make([]v1.MemberInfo, 0, len(members))
	for _, m := range members {
		info := v1.MemberInfo{
			UserID:       m.UserID,
			DepartmentID: m.DepartmentID,
			PositionID:   m.PositionID,
			IsLeader:     m.IsLeader,
			JoinedAt:     m.JoinedAt,
		}
		if m.PositionID > 0 {
			info.PositionName = posNames[m.PositionID]
		}
		result = append(result, info)
	}
	return result, total, nil
}

// departmentNameMap 批量查部门名，返回 id -> name 的映射（缺失的 id 为空串）。
func departmentNameMap(ids []uint) map[uint]string {
	m := make(map[uint]string, len(ids))
	if len(ids) == 0 {
		return m
	}
	depts, err := dao.GetDepartmentDAO().GetByIDs(ids)
	if err != nil {
		return m
	}
	for _, d := range depts {
		m[d.ID] = d.Name
	}
	return m
}

// positionNameMap 批量查岗位名，返回 id -> name 的映射。
func positionNameMap(ids []uint) map[uint]string {
	m := make(map[uint]string, len(ids))
	if len(ids) == 0 {
		return m
	}
	positions, err := dao.GetPositionDAO().GetByIDs(ids)
	if err != nil {
		return m
	}
	for _, p := range positions {
		m[p.ID] = p.Name
	}
	return m
}

func (s *memberServiceImpl) GetMember(userID, departmentID uint) (*model.OrgMember, error) {
	return dao.GetMemberDAO().GetByUserAndDept(userID, departmentID)
}
