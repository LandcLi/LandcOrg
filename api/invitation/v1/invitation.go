package v1

import "time"

// --- Invitation API ---

// CreateInvitationRequest creates an invitation to join a department.
type CreateInvitationRequest struct {
	InviterID    uint `json:"inviter_id" binding:"required"`
	DepartmentID uint `json:"department_id" binding:"required"`
	PositionID   uint `json:"position_id"`
	InviteeID    uint `json:"invitee_id"`
	MaxUses      int  `json:"max_uses"`
	ExpiresInHours int `json:"expires_in_hours"`
}

type CreateInvitationResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AcceptInvitationRequest accepts an invitation by code.
type AcceptInvitationRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type AcceptInvitationResponse struct {
	Success      bool `json:"success"`
	DepartmentID uint `json:"department_id"`
}

// RevokeInvitationRequest revokes an invitation.
type RevokeInvitationRequest struct {
	InviterID    uint `json:"inviter_id" binding:"required"`
	InvitationID uint `json:"invitation_id" binding:"required"`
}

type RevokeInvitationResponse struct {
	Success bool `json:"success"`
}

// ListInvitationsRequest lists invitations created by a user.
type ListInvitationsRequest struct {
	InviterID uint `json:"inviter_id" binding:"required"`
	Page      int  `json:"page"`
	PageSize  int  `json:"page_size"`
}

type InvitationInfo struct {
	ID             uint      `json:"id"`
	Code           string    `json:"code"`
	DepartmentID   uint      `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	PositionID     uint      `json:"position_id"`
	InviteeID      uint      `json:"invitee_id"`
	Status         int       `json:"status"`
	MaxUses        int       `json:"max_uses"`
	UsedCount      int       `json:"used_count"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type ListInvitationsResponse struct {
	Total       int64            `json:"total"`
	Invitations []InvitationInfo `json:"invitations"`
}
