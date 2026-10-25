package v1

import (
	"github.com/LandcLi/landc-go/frame/pkg/meta"
)

import "time"

// AddMemberRequest adds a user to a department.
type AddMemberRequest struct {
	meta.Meta    `path:"/api/member/add" method:"POST"`
	UserID       uint `json:"user_id" binding:"required"`
	DepartmentID uint `json:"department_id" binding:"required"`
	PositionID   uint `json:"position_id"`
	IsLeader     bool `json:"is_leader"`
}

type AddMemberResponse struct {
	Success bool `json:"success"`
}

// RemoveMemberRequest removes a user from a department.
type RemoveMemberRequest struct {
	meta.Meta    `path:"/api/member/remove" method:"POST"`
	UserID       uint `json:"user_id" binding:"required"`
	DepartmentID uint `json:"department_id" binding:"required"`
}

type RemoveMemberResponse struct {
	Success bool `json:"success"`
}

// UpdateMemberRequest updates a member's position or leader status.
type UpdateMemberRequest struct {
	meta.Meta    `path:"/api/member/update" method:"POST"`
	UserID       uint  `json:"user_id" binding:"required"`
	DepartmentID uint  `json:"department_id" binding:"required"`
	PositionID   *uint `json:"position_id"`
	IsLeader     *bool `json:"is_leader"`
}

type UpdateMemberResponse struct {
	Success bool `json:"success"`
}

// GetUserDepartmentsRequest gets all departments a user belongs to.
type GetUserDepartmentsRequest struct {
	meta.Meta `path:"/api/member/get-user-departments" method:"POST"`
	UserID    uint `json:"user_id" binding:"required"`
}

// MemberInfo represents a member's organization info.
type MemberInfo struct {
	UserID         uint      `json:"user_id"`
	DepartmentID   uint      `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	PositionID     uint      `json:"position_id"`
	PositionName   string    `json:"position_name"`
	IsLeader       bool      `json:"is_leader"`
	JoinedAt       time.Time `json:"joined_at"`
}

type GetUserDepartmentsResponse struct {
	Memberships []MemberInfo `json:"memberships"`
}

// GetDepartmentMembersRequest gets all members of a department.
type GetDepartmentMembersRequest struct {
	meta.Meta    `path:"/api/member/get-department-members" method:"POST"`
	DepartmentID uint `json:"department_id" binding:"required"`
	Page         int  `json:"page"`
	PageSize     int  `json:"page_size"`
}

type GetDepartmentMembersResponse struct {
	Total   int64        `json:"total"`
	Members []MemberInfo `json:"members"`
}
