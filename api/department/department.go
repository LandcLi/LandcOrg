package department

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

// DepartmentController defines the department management interface.
type DepartmentController interface {
	CreateDepartment(ctx context.Context, req *v1.CreateDepartmentRequest) (*v1.CreateDepartmentResponse, error)
	UpdateDepartment(ctx context.Context, req *v1.UpdateDepartmentRequest) (*v1.UpdateDepartmentResponse, error)
	DeleteDepartment(ctx context.Context, req *v1.DeleteDepartmentRequest) (*v1.DeleteDepartmentResponse, error)
	GetDepartment(ctx context.Context, req *v1.GetDepartmentRequest) (*v1.GetDepartmentResponse, error)
	GetDepartmentTree(ctx context.Context, req *v1.GetDepartmentTreeRequest) (*v1.GetDepartmentTreeResponse, error)
	ListDepartments(ctx context.Context, req *v1.ListDepartmentsRequest) (*v1.ListDepartmentsResponse, error)
}

var DepartmentRoutes = di.Routes{
	"CreateDepartment": {Method: "POST", Path: "/api/department/create"},
	"UpdateDepartment": {Method: "POST", Path: "/api/department/update"},
	"DeleteDepartment": {Method: "POST", Path: "/api/department/delete"},
	"GetDepartment":    {Method: "POST", Path: "/api/department/get"},
	"GetDepartmentTree": {Method: "POST", Path: "/api/department/tree"},
	"ListDepartments":  {Method: "POST", Path: "/api/department/list"},
}

var DepartmentGateway = di.NewGateway[DepartmentController]("department.controller", DepartmentRoutes)

func GetDepartmentController() DepartmentController {
	return DepartmentGateway.Get()
}
