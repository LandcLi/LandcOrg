package member

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/member/v1"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// MemberController defines the organization member management interface.
type MemberController interface {
	AddMember(ctx context.Context, req *v1.AddMemberRequest) (*v1.AddMemberResponse, error)
	RemoveMember(ctx context.Context, req *v1.RemoveMemberRequest) (*v1.RemoveMemberResponse, error)
	UpdateMember(ctx context.Context, req *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error)
	GetUserDepartments(ctx context.Context, req *v1.GetUserDepartmentsRequest) (*v1.GetUserDepartmentsResponse, error)
	GetDepartmentMembers(ctx context.Context, req *v1.GetDepartmentMembersRequest) (*v1.GetDepartmentMembersResponse, error)
}

var MemberGateway = di.NewGateway[MemberController]("member.controller")

func GetMemberController() MemberController {
	return MemberGateway.Get()
}
