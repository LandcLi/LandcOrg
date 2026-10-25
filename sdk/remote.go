package sdk

import (
	"context"

	deptV1 "github.com/LandcLi/LandcOrg/api/department/v1"
	invV1 "github.com/LandcLi/LandcOrg/api/invitation/v1"
	memV1 "github.com/LandcLi/LandcOrg/api/member/v1"
	orgRoleV1 "github.com/LandcLi/LandcOrg/api/org_role/v1"
	posV1 "github.com/LandcLi/LandcOrg/api/position/v1"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// --- Remote Department Controller ---

type RemoteDepartmentController struct{ client *di.ProxyClient }

func NewRemoteDepartmentController(baseURL string, opts ...di.RemoteOption) *RemoteDepartmentController {
	return &RemoteDepartmentController{client: di.NewProxyClient(baseURL, opts...)}
}

func (r *RemoteDepartmentController) CreateDepartment(ctx context.Context, req *deptV1.CreateDepartmentRequest) (*deptV1.CreateDepartmentResponse, error) {
	return di.Call[deptV1.CreateDepartmentResponse](r.client, ctx, "CreateDepartment", req)
}
func (r *RemoteDepartmentController) UpdateDepartment(ctx context.Context, req *deptV1.UpdateDepartmentRequest) (*deptV1.UpdateDepartmentResponse, error) {
	return di.Call[deptV1.UpdateDepartmentResponse](r.client, ctx, "UpdateDepartment", req)
}
func (r *RemoteDepartmentController) DeleteDepartment(ctx context.Context, req *deptV1.DeleteDepartmentRequest) (*deptV1.DeleteDepartmentResponse, error) {
	return di.Call[deptV1.DeleteDepartmentResponse](r.client, ctx, "DeleteDepartment", req)
}
func (r *RemoteDepartmentController) GetDepartment(ctx context.Context, req *deptV1.GetDepartmentRequest) (*deptV1.GetDepartmentResponse, error) {
	return di.Call[deptV1.GetDepartmentResponse](r.client, ctx, "GetDepartment", req)
}
func (r *RemoteDepartmentController) GetDepartmentTree(ctx context.Context, req *deptV1.GetDepartmentTreeRequest) (*deptV1.GetDepartmentTreeResponse, error) {
	return di.Call[deptV1.GetDepartmentTreeResponse](r.client, ctx, "GetDepartmentTree", req)
}
func (r *RemoteDepartmentController) ListDepartments(ctx context.Context, req *deptV1.ListDepartmentsRequest) (*deptV1.ListDepartmentsResponse, error) {
	return di.Call[deptV1.ListDepartmentsResponse](r.client, ctx, "ListDepartments", req)
}

// --- Remote Position Controller ---

type RemotePositionController struct{ client *di.ProxyClient }

func NewRemotePositionController(baseURL string, opts ...di.RemoteOption) *RemotePositionController {
	return &RemotePositionController{client: di.NewProxyClient(baseURL, opts...)}
}

func (r *RemotePositionController) CreatePosition(ctx context.Context, req *posV1.CreatePositionRequest) (*posV1.CreatePositionResponse, error) {
	return di.Call[posV1.CreatePositionResponse](r.client, ctx, "CreatePosition", req)
}
func (r *RemotePositionController) UpdatePosition(ctx context.Context, req *posV1.UpdatePositionRequest) (*posV1.UpdatePositionResponse, error) {
	return di.Call[posV1.UpdatePositionResponse](r.client, ctx, "UpdatePosition", req)
}
func (r *RemotePositionController) DeletePosition(ctx context.Context, req *posV1.DeletePositionRequest) (*posV1.DeletePositionResponse, error) {
	return di.Call[posV1.DeletePositionResponse](r.client, ctx, "DeletePosition", req)
}
func (r *RemotePositionController) GetPosition(ctx context.Context, req *posV1.GetPositionRequest) (*posV1.GetPositionResponse, error) {
	return di.Call[posV1.GetPositionResponse](r.client, ctx, "GetPosition", req)
}
func (r *RemotePositionController) ListPositions(ctx context.Context, req *posV1.ListPositionsRequest) (*posV1.ListPositionsResponse, error) {
	return di.Call[posV1.ListPositionsResponse](r.client, ctx, "ListPositions", req)
}

// --- Remote Member Controller ---

type RemoteMemberController struct{ client *di.ProxyClient }

func NewRemoteMemberController(baseURL string, opts ...di.RemoteOption) *RemoteMemberController {
	return &RemoteMemberController{client: di.NewProxyClient(baseURL, opts...)}
}

func (r *RemoteMemberController) AddMember(ctx context.Context, req *memV1.AddMemberRequest) (*memV1.AddMemberResponse, error) {
	return di.Call[memV1.AddMemberResponse](r.client, ctx, "AddMember", req)
}
func (r *RemoteMemberController) RemoveMember(ctx context.Context, req *memV1.RemoveMemberRequest) (*memV1.RemoveMemberResponse, error) {
	return di.Call[memV1.RemoveMemberResponse](r.client, ctx, "RemoveMember", req)
}
func (r *RemoteMemberController) UpdateMember(ctx context.Context, req *memV1.UpdateMemberRequest) (*memV1.UpdateMemberResponse, error) {
	return di.Call[memV1.UpdateMemberResponse](r.client, ctx, "UpdateMember", req)
}
func (r *RemoteMemberController) GetUserDepartments(ctx context.Context, req *memV1.GetUserDepartmentsRequest) (*memV1.GetUserDepartmentsResponse, error) {
	return di.Call[memV1.GetUserDepartmentsResponse](r.client, ctx, "GetUserDepartments", req)
}
func (r *RemoteMemberController) GetDepartmentMembers(ctx context.Context, req *memV1.GetDepartmentMembersRequest) (*memV1.GetDepartmentMembersResponse, error) {
	return di.Call[memV1.GetDepartmentMembersResponse](r.client, ctx, "GetDepartmentMembers", req)
}

// --- Remote Invitation Controller ---

type RemoteInvitationController struct{ client *di.ProxyClient }

func NewRemoteInvitationController(baseURL string, opts ...di.RemoteOption) *RemoteInvitationController {
	return &RemoteInvitationController{client: di.NewProxyClient(baseURL, opts...)}
}

func (r *RemoteInvitationController) CreateInvitation(ctx context.Context, req *invV1.CreateInvitationRequest) (*invV1.CreateInvitationResponse, error) {
	return di.Call[invV1.CreateInvitationResponse](r.client, ctx, "CreateInvitation", req)
}
func (r *RemoteInvitationController) AcceptInvitation(ctx context.Context, req *invV1.AcceptInvitationRequest) (*invV1.AcceptInvitationResponse, error) {
	return di.Call[invV1.AcceptInvitationResponse](r.client, ctx, "AcceptInvitation", req)
}
func (r *RemoteInvitationController) RevokeInvitation(ctx context.Context, req *invV1.RevokeInvitationRequest) (*invV1.RevokeInvitationResponse, error) {
	return di.Call[invV1.RevokeInvitationResponse](r.client, ctx, "RevokeInvitation", req)
}
func (r *RemoteInvitationController) ListInvitations(ctx context.Context, req *invV1.ListInvitationsRequest) (*invV1.ListInvitationsResponse, error) {
	return di.Call[invV1.ListInvitationsResponse](r.client, ctx, "ListInvitations", req)
}

// --- Remote OrgRole Controller ---

type RemoteOrgRoleController struct{ client *di.ProxyClient }

func NewRemoteOrgRoleController(baseURL string, opts ...di.RemoteOption) *RemoteOrgRoleController {
	return &RemoteOrgRoleController{client: di.NewProxyClient(baseURL, opts...)}
}

func (r *RemoteOrgRoleController) AssignDepartmentRoles(ctx context.Context, req *orgRoleV1.AssignDepartmentRolesRequest) (*orgRoleV1.AssignDepartmentRolesResponse, error) {
	return di.Call[orgRoleV1.AssignDepartmentRolesResponse](r.client, ctx, "AssignDepartmentRoles", req)
}
func (r *RemoteOrgRoleController) RevokeDepartmentRoles(ctx context.Context, req *orgRoleV1.RevokeDepartmentRolesRequest) (*orgRoleV1.RevokeDepartmentRolesResponse, error) {
	return di.Call[orgRoleV1.RevokeDepartmentRolesResponse](r.client, ctx, "RevokeDepartmentRoles", req)
}
func (r *RemoteOrgRoleController) GetDepartmentRoles(ctx context.Context, req *orgRoleV1.GetDepartmentRolesRequest) (*orgRoleV1.GetDepartmentRolesResponse, error) {
	return di.Call[orgRoleV1.GetDepartmentRolesResponse](r.client, ctx, "GetDepartmentRoles", req)
}
func (r *RemoteOrgRoleController) AssignPositionRoles(ctx context.Context, req *orgRoleV1.AssignPositionRolesRequest) (*orgRoleV1.AssignPositionRolesResponse, error) {
	return di.Call[orgRoleV1.AssignPositionRolesResponse](r.client, ctx, "AssignPositionRoles", req)
}
func (r *RemoteOrgRoleController) RevokePositionRoles(ctx context.Context, req *orgRoleV1.RevokePositionRolesRequest) (*orgRoleV1.RevokePositionRolesResponse, error) {
	return di.Call[orgRoleV1.RevokePositionRolesResponse](r.client, ctx, "RevokePositionRoles", req)
}
func (r *RemoteOrgRoleController) GetPositionRoles(ctx context.Context, req *orgRoleV1.GetPositionRolesRequest) (*orgRoleV1.GetPositionRolesResponse, error) {
	return di.Call[orgRoleV1.GetPositionRolesResponse](r.client, ctx, "GetPositionRoles", req)
}
func (r *RemoteOrgRoleController) GetUserInheritedRoles(ctx context.Context, req *orgRoleV1.GetUserInheritedRolesRequest) (*orgRoleV1.GetUserInheritedRolesResponse, error) {
	return di.Call[orgRoleV1.GetUserInheritedRolesResponse](r.client, ctx, "GetUserInheritedRoles", req)
}
