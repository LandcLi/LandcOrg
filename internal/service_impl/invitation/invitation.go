package invitation

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterInvitationService(&invitationServiceImpl{}) }

type invitationServiceImpl struct{}

func (s *invitationServiceImpl) CreateInvitation(inviterID, departmentID, positionID uint, inviteeID uint, maxUses int, expiresInHours int) (*model.OrgInvitation, error) {
	// Verify department exists
	dept, err := dao.GetDepartmentDAO().GetByID(departmentID)
	if err != nil || dept == nil {
		return nil, fmt.Errorf("department not found")
	}

	if maxUses == 0 {
		maxUses = 1
	}
	if expiresInHours <= 0 {
		expiresInHours = 72 // Default 3 days
	}

	code := generateInviteCode()

	inv := &model.OrgInvitation{
		Code:         code,
		DepartmentID: departmentID,
		PositionID:   positionID,
		InviterID:    inviterID,
		InviteeID:    inviteeID,
		Status:       0,
		MaxUses:      maxUses,
		UsedCount:    0,
		ExpiresAt:    time.Now().Add(time.Duration(expiresInHours) * time.Hour),
	}

	if err := dao.GetInvitationDAO().Create(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *invitationServiceImpl) AcceptInvitation(userID uint, code string) error {
	inv, err := dao.GetInvitationDAO().GetByCode(code)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	// Validate invitation
	if inv.Status != 0 {
		return fmt.Errorf("invitation is no longer valid (status=%d)", inv.Status)
	}
	if time.Now().After(inv.ExpiresAt) {
		inv.Status = 2 // Mark as expired
		_ = dao.GetInvitationDAO().Update(inv)
		return fmt.Errorf("invitation has expired")
	}
	if inv.MaxUses > 0 && inv.UsedCount >= inv.MaxUses {
		return fmt.Errorf("invitation has reached max uses")
	}
	if inv.InviteeID > 0 && inv.InviteeID != userID {
		return fmt.Errorf("this invitation is not for you")
	}

	// Add user to department
	err = service.GetMemberService().AddMember(userID, inv.DepartmentID, inv.PositionID, false)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	// Update invitation usage
	inv.UsedCount++
	if inv.MaxUses > 0 && inv.UsedCount >= inv.MaxUses {
		inv.Status = 1 // Mark as fully used
	}
	_ = dao.GetInvitationDAO().Update(inv)

	return nil
}

func (s *invitationServiceImpl) RevokeInvitation(inviterID uint, invitationID uint) error {
	inv, err := dao.GetInvitationDAO().GetByID(invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}
	if inv.InviterID != inviterID {
		return fmt.Errorf("you can only revoke your own invitations")
	}
	if inv.Status != 0 {
		return fmt.Errorf("invitation is already in terminal state")
	}
	inv.Status = 3 // Revoked
	return dao.GetInvitationDAO().Update(inv)
}

func (s *invitationServiceImpl) ListMyInvitations(inviterID uint, page, pageSize int) ([]model.OrgInvitation, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return dao.GetInvitationDAO().ListByInviter(inviterID, page, pageSize)
}

func generateInviteCode() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
