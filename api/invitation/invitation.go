package invitation

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/invitation/v1"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// InvitationController defines the invitation management interface.
type InvitationController interface {
	CreateInvitation(ctx context.Context, req *v1.CreateInvitationRequest) (*v1.CreateInvitationResponse, error)
	AcceptInvitation(ctx context.Context, req *v1.AcceptInvitationRequest) (*v1.AcceptInvitationResponse, error)
	RevokeInvitation(ctx context.Context, req *v1.RevokeInvitationRequest) (*v1.RevokeInvitationResponse, error)
	ListInvitations(ctx context.Context, req *v1.ListInvitationsRequest) (*v1.ListInvitationsResponse, error)
}

var InvitationGateway = di.NewGateway[InvitationController]("invitation.controller")

func GetInvitationController() InvitationController {
	return InvitationGateway.Get()
}
