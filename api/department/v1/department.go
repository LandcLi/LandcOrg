package v1

import (
	"github.com/LandcLi/landc-go/frame/pkg/meta"
)

// CreateDepartmentRequest creates a new department.
type CreateDepartmentRequest struct {
	meta.Meta   `path:"/api/department/create" method:"POST"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	Leader      string `json:"leader"`
	Description string `json:"description"`
}

type CreateDepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateDepartmentRequest struct {
	meta.Meta   `path:"/api/department/update" method:"POST"`
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	ParentID    *uint  `json:"parent_id"`
	Sort        *int   `json:"sort"`
	Leader      string `json:"leader"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

type UpdateDepartmentResponse struct {
	Success bool `json:"success"`
}

type DeleteDepartmentRequest struct {
	meta.Meta `path:"/api/department/delete" method:"POST"`
	ID        uint `json:"id" binding:"required"`
}

type DeleteDepartmentResponse struct {
	Success bool `json:"success"`
}

type GetDepartmentRequest struct {
	meta.Meta `path:"/api/department/get" method:"POST"`
	ID        uint `json:"id" binding:"required"`
}

type GetDepartmentResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	ParentID    uint   `json:"parent_id"`
	Path        string `json:"path"`
	Level       int    `json:"level"`
	Sort        int    `json:"sort"`
	Leader      string `json:"leader"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// GetDepartmentTreeRequest gets the full department tree.
type GetDepartmentTreeRequest struct {
	meta.Meta `path:"/api/department/tree" method:"POST"`
	RootID    uint `json:"root_id"` // 0 means get all from root
}

// DepartmentTreeNode represents a node in the department tree.
type DepartmentTreeNode struct {
	ID       uint                  `json:"id"`
	Name     string                `json:"name"`
	Code     string                `json:"code"`
	ParentID uint                  `json:"parent_id"`
	Level    int                   `json:"level"`
	Sort     int                   `json:"sort"`
	Leader   string                `json:"leader"`
	Status   int                   `json:"status"`
	Children []*DepartmentTreeNode `json:"children"`
}

type GetDepartmentTreeResponse struct {
	Tree []*DepartmentTreeNode `json:"tree"`
}

// ListDepartmentsRequest lists departments with pagination.
type ListDepartmentsRequest struct {
	meta.Meta `path:"/api/department/list" method:"POST"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	ParentID  *uint  `json:"parent_id"`
	Keyword   string `json:"keyword"`
}

type ListDepartmentsResponse struct {
	Total int64                   `json:"total"`
	List  []GetDepartmentResponse `json:"list"`
}
