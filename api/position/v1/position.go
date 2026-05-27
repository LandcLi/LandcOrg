package v1

// CreatePositionRequest creates a new position.
type CreatePositionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type CreatePositionResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdatePositionRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Level       *int   `json:"level"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

type UpdatePositionResponse struct {
	Success bool `json:"success"`
}

type DeletePositionRequest struct {
	ID uint `json:"id" binding:"required"`
}

type DeletePositionResponse struct {
	Success bool `json:"success"`
}

type GetPositionRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetPositionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Level       int    `json:"level"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type ListPositionsRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type ListPositionsResponse struct {
	Total int64                `json:"total"`
	List  []GetPositionResponse `json:"list"`
}
