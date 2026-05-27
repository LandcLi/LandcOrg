package v1

// CreateDepartmentRequest creates a new department.
type CreateDepartmentRequest struct {
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
	ID uint `json:"id" binding:"required"`
}

type DeleteDepartmentResponse struct {
	Success bool `json:"success"`
}

type GetDepartmentRequest struct {
	ID uint `json:"id" binding:"required"`
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
	RootID uint `json:"root_id"` // 0 means get all from root
}

// DepartmentTreeNode represents a node in the department tree.
type DepartmentTreeNode struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Code        string                `json:"code"`
	ParentID    uint                  `json:"parent_id"`
	Level       int                   `json:"level"`
	Sort        int                   `json:"sort"`
	Leader      string                `json:"leader"`
	Status      int                   `json:"status"`
	Children    []*DepartmentTreeNode `json:"children"`
}

type GetDepartmentTreeResponse struct {
	Tree []*DepartmentTreeNode `json:"tree"`
}

// ListDepartmentsRequest lists departments with pagination.
type ListDepartmentsRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	ParentID *uint  `json:"parent_id"`
	Keyword  string `json:"keyword"`
}

type ListDepartmentsResponse struct {
	Total int64                  `json:"total"`
	List  []GetDepartmentResponse `json:"list"`
}
