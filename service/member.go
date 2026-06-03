package service

import (
	"github.com/LandcLi/LandcOrg/api/member/v1"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// MemberService defines the organization member business logic interface.
type MemberService interface {
	AddMember(userID, departmentID, positionID uint, isLeader bool) error
	RemoveMember(userID, departmentID uint) error
	UpdateMember(userID, departmentID uint, positionID *uint, isLeader *bool) error
	GetUserDepartments(userID uint) ([]v1.MemberInfo, error)
	GetDepartmentMembers(departmentID uint, page, pageSize int) ([]v1.MemberInfo, int64, error)
	GetMember(userID, departmentID uint) (*model.OrgMember, error)
}

func RegisterMemberService(impl MemberService) {
	di.Provide[MemberService]("member.service", impl)
}

func GetMemberService() MemberService {
	return di.Require[MemberService]("member.service")
}
