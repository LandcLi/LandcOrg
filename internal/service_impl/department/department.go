package department

import (
	"fmt"

	v1 "github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterDepartmentService(&deptServiceImpl{}) }

type deptServiceImpl struct{}

func (s *deptServiceImpl) CreateDepartment(name, code string, parentID uint, sort int, leader, description string) (*model.Department, error) {
	existing, _ := dao.GetDepartmentDAO().GetByCode(code)
	if existing != nil && existing.ID != 0 {
		return nil, fmt.Errorf("department code '%s' already exists", code)
	}

	dept := &model.Department{
		Name:        name,
		Code:        code,
		ParentID:    parentID,
		Sort:        sort,
		Leader:      leader,
		Description: description,
		Status:      1,
		Level:       1,
	}

	// Calculate level and path
	if parentID > 0 {
		parent, err := dao.GetDepartmentDAO().GetByID(parentID)
		if err != nil {
			return nil, fmt.Errorf("parent department not found: %w", err)
		}
		dept.Level = parent.Level + 1
		dept.Path = fmt.Sprintf("%s%d/", parent.Path, parent.ID)
	} else {
		dept.Path = "/"
	}

	if err := dao.GetDepartmentDAO().Create(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *deptServiceImpl) UpdateDepartment(id uint, name, code string, parentID *uint, sort *int, leader, description string, status *int) error {
	dept, err := dao.GetDepartmentDAO().GetByID(id)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}
	if name != "" {
		dept.Name = name
	}
	if code != "" {
		existing, _ := dao.GetDepartmentDAO().GetByCode(code)
		if existing != nil && existing.ID != id {
			return fmt.Errorf("department code '%s' already exists", code)
		}
		dept.Code = code
	}
	if parentID != nil {
		dept.ParentID = *parentID
	}
	if sort != nil {
		dept.Sort = *sort
	}
	if leader != "" {
		dept.Leader = leader
	}
	if description != "" {
		dept.Description = description
	}
	if status != nil {
		dept.Status = *status
	}
	return dao.GetDepartmentDAO().Update(dept)
}

func (s *deptServiceImpl) DeleteDepartment(id uint) error {
	// Check if has children
	children, _ := dao.GetDepartmentDAO().GetByParentID(id)
	if len(children) > 0 {
		return fmt.Errorf("cannot delete department with sub-departments")
	}
	return dao.GetDepartmentDAO().Delete(id)
}

func (s *deptServiceImpl) GetDepartment(id uint) (*model.Department, error) {
	return dao.GetDepartmentDAO().GetByID(id)
}

func (s *deptServiceImpl) GetDepartmentTree(rootID uint) ([]*v1.DepartmentTreeNode, error) {
	allDepts, err := dao.GetDepartmentDAO().GetAll()
	if err != nil {
		return nil, err
	}

	// Build tree
	nodeMap := make(map[uint]*v1.DepartmentTreeNode)
	for _, d := range allDepts {
		nodeMap[d.ID] = &v1.DepartmentTreeNode{
			ID:       d.ID,
			Name:     d.Name,
			Code:     d.Code,
			ParentID: d.ParentID,
			Level:    d.Level,
			Sort:     d.Sort,
			Leader:   d.Leader,
			Status:   d.Status,
			Children: make([]*v1.DepartmentTreeNode, 0),
		}
	}

	var roots []*v1.DepartmentTreeNode
	for _, d := range allDepts {
		node := nodeMap[d.ID]
		if d.ParentID == rootID {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[d.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

func (s *deptServiceImpl) ListDepartments(page, pageSize int, parentID *uint, keyword string) ([]model.Department, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return dao.GetDepartmentDAO().List(page, pageSize, parentID, keyword)
}
