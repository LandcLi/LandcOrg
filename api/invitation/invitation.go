package invitation

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/invitation/v1"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

// InvitationController defines the invitation management interface.
type InvitationController interface {
	CreateInvitation(ctx context.Context, req *v1.CreateInvitationRequest) (*v1.CreateInvitationResponse, error)
	AcceptInvitation(ctx context.Context, req *v1.AcceptInvitationRequest) (*v1.AcceptInvitationResponse, error)
	RevokeInvitation(ctx context.Context, req *v1.RevokeInvitationRequest) (*v1.RevokeInvitationResponse, error)
	ListInvitations(ctx context.Context, req *v1.ListInvitationsRequest) (*v1.ListInvitationsResponse, error)
}

var InvitationRoutes = di.Routes{
	"CreateInvitation": {Method: "POST", Path: "/api/invitation/create"},
	"AcceptInvitation": {Method: "POST", Path: "/api/invitation/accept"},
	"RevokeInvitation": {Method: "POST", Path: "/api/invitation/revoke"},
	"ListInvitations":  {Method: "POST", Path: "/api/invitation/list"},
}

var InvitationGateway = di.NewGateway[InvitationController]("invitation.controller", InvitationRoutes)

func GetInvitationController() InvitationController {
	return InvitationGateway.Get()
}
