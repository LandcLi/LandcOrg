package sdk

import (
	"context"

	"github.com/LandcLi/LandcOrg/api/department"
	deptV1 "github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/LandcOrg/api/member"
	memV1 "github.com/LandcLi/LandcOrg/api/member/v1"
	"github.com/LandcLi/LandcOrg/api/position"
	posV1 "github.com/LandcLi/LandcOrg/api/position/v1"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

// --- Remote Department Controller ---

type RemoteDepartmentController struct{ client *di.ProxyClient }

func NewRemoteDepartmentController(baseURL string, opts ...di.RemoteOption) *RemoteDepartmentController {
	return &RemoteDepartmentController{client: di.NewProxyClient(baseURL, department.DepartmentRoutes, opts...)}
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
	return &RemotePositionController{client: di.NewProxyClient(baseURL, position.PositionRoutes, opts...)}
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
	return &RemoteMemberController{client: di.NewProxyClient(baseURL, member.MemberRoutes, opts...)}
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
