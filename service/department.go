package service

import (
	"github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// DepartmentService defines the department business logic interface.
type DepartmentService interface {
	CreateDepartment(name, code string, parentID uint, sort int, leader, description string) (*model.Department, error)
	UpdateDepartment(id uint, name, code string, parentID *uint, sort *int, leader, description string, status *int) error
	DeleteDepartment(id uint) error
	GetDepartment(id uint) (*model.Department, error)
	GetDepartmentTree(rootID uint) ([]*v1.DepartmentTreeNode, error)
	ListDepartments(page, pageSize int, parentID *uint, keyword string) ([]model.Department, int64, error)
}

func RegisterDepartmentService(impl DepartmentService) {
	di.Provide[DepartmentService]("department.service", impl)
}

func GetDepartmentService() DepartmentService {
	return di.Require[DepartmentService]("department.service")
}
