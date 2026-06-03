package service

import (
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// InvitationService defines the invitation business logic interface.
type InvitationService interface {
	// CreateInvitation creates an invitation for a department.
	CreateInvitation(inviterID, departmentID, positionID uint, inviteeID uint, maxUses int, expiresInHours int) (*model.OrgInvitation, error)
	// AcceptInvitation accepts an invitation by code.
	AcceptInvitation(userID uint, code string) error
	// RevokeInvitation revokes an invitation.
	RevokeInvitation(inviterID uint, invitationID uint) error
	// ListMyInvitations lists invitations created by a user.
	ListMyInvitations(inviterID uint, page, pageSize int) ([]model.OrgInvitation, int64, error)
}

func RegisterInvitationService(impl InvitationService) {
	di.Provide[InvitationService]("invitation.service", impl)
}

func GetInvitationService() InvitationService {
	return di.Require[InvitationService]("invitation.service")
}
