package department

import (
	"context"

	deptApi "github.com/LandcLi/LandcOrg/api/department"
	v1 "github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/LandcOrg/service"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

func init() { di.Provide[deptApi.DepartmentController]("department.controller", &deptCtrl{}) }

type deptCtrl struct{}

func (c *deptCtrl) CreateDepartment(ctx context.Context, req *v1.CreateDepartmentRequest) (*v1.CreateDepartmentResponse, error) {
	dept, err := service.GetDepartmentService().CreateDepartment(req.Name, req.Code, req.ParentID, req.Sort, req.Leader, req.Description)
	if err != nil {
		return nil, err
	}
	return &v1.CreateDepartmentResponse{ID: dept.ID, Name: dept.Name, Code: dept.Code}, nil
}

func (c *deptCtrl) UpdateDepartment(ctx context.Context, req *v1.UpdateDepartmentRequest) (*v1.UpdateDepartmentResponse, error) {
	err := service.GetDepartmentService().UpdateDepartment(req.ID, req.Name, req.Code, req.ParentID, req.Sort, req.Leader, req.Description, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateDepartmentResponse{Success: true}, nil
}

func (c *deptCtrl) DeleteDepartment(ctx context.Context, req *v1.DeleteDepartmentRequest) (*v1.DeleteDepartmentResponse, error) {
	if err := service.GetDepartmentService().DeleteDepartment(req.ID); err != nil {
		return nil, err
	}
	return &v1.DeleteDepartmentResponse{Success: true}, nil
}

func (c *deptCtrl) GetDepartment(ctx context.Context, req *v1.GetDepartmentRequest) (*v1.GetDepartmentResponse, error) {
	dept, err := service.GetDepartmentService().GetDepartment(req.ID)
	if err != nil {
		return nil, err
	}
	return &v1.GetDepartmentResponse{
		ID: dept.ID, Name: dept.Name, Code: dept.Code, ParentID: dept.ParentID,
		Path: dept.Path, Level: dept.Level, Sort: dept.Sort, Leader: dept.Leader,
		Description: dept.Description, Status: dept.Status,
	}, nil
}

func (c *deptCtrl) GetDepartmentTree(ctx context.Context, req *v1.GetDepartmentTreeRequest) (*v1.GetDepartmentTreeResponse, error) {
	tree, err := service.GetDepartmentService().GetDepartmentTree(req.RootID)
	if err != nil {
		return nil, err
	}
	return &v1.GetDepartmentTreeResponse{Tree: tree}, nil
}

func (c *deptCtrl) ListDepartments(ctx context.Context, req *v1.ListDepartmentsRequest) (*v1.ListDepartmentsResponse, error) {
	depts, total, err := service.GetDepartmentService().ListDepartments(req.Page, req.PageSize, req.ParentID, req.Keyword)
	if err != nil {
		return nil, err
	}
	list := make([]v1.GetDepartmentResponse, 0, len(depts))
	for _, d := range depts {
		list = append(list, v1.GetDepartmentResponse{
			ID: d.ID, Name: d.Name, Code: d.Code, ParentID: d.ParentID,
			Path: d.Path, Level: d.Level, Sort: d.Sort, Leader: d.Leader,
			Description: d.Description, Status: d.Status,
		})
	}
	return &v1.ListDepartmentsResponse{Total: total, List: list}, nil
}
