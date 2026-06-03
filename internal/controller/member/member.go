package member

import (
	"context"

	memApi "github.com/LandcLi/LandcOrg/api/member"
	v1 "github.com/LandcLi/LandcOrg/api/member/v1"
	"github.com/LandcLi/LandcOrg/service"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

func init() { di.Provide[memApi.MemberController]("member.controller", &memCtrl{}) }

type memCtrl struct{}

func (c *memCtrl) AddMember(ctx context.Context, req *v1.AddMemberRequest) (*v1.AddMemberResponse, error) {
	if err := service.GetMemberService().AddMember(req.UserID, req.DepartmentID, req.PositionID, req.IsLeader); err != nil {
		return nil, err
	}
	return &v1.AddMemberResponse{Success: true}, nil
}

func (c *memCtrl) RemoveMember(ctx context.Context, req *v1.RemoveMemberRequest) (*v1.RemoveMemberResponse, error) {
	if err := service.GetMemberService().RemoveMember(req.UserID, req.DepartmentID); err != nil {
		return nil, err
	}
	return &v1.RemoveMemberResponse{Success: true}, nil
}

func (c *memCtrl) UpdateMember(ctx context.Context, req *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	if err := service.GetMemberService().UpdateMember(req.UserID, req.DepartmentID, req.PositionID, req.IsLeader); err != nil {
		return nil, err
	}
	return &v1.UpdateMemberResponse{Success: true}, nil
}

func (c *memCtrl) GetUserDepartments(ctx context.Context, req *v1.GetUserDepartmentsRequest) (*v1.GetUserDepartmentsResponse, error) {
	memberships, err := service.GetMemberService().GetUserDepartments(req.UserID)
	if err != nil { return nil, err }
	return &v1.GetUserDepartmentsResponse{Memberships: memberships}, nil
}

func (c *memCtrl) GetDepartmentMembers(ctx context.Context, req *v1.GetDepartmentMembersRequest) (*v1.GetDepartmentMembersResponse, error) {
	members, total, err := service.GetMemberService().GetDepartmentMembers(req.DepartmentID, req.Page, req.PageSize)
	if err != nil { return nil, err }
	return &v1.GetDepartmentMembersResponse{Total: total, Members: members}, nil
}
