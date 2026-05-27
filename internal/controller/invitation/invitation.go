package invitation

import (
	"context"

	invApi "github.com/LandcLi/LandcOrg/api/invitation"
	v1 "github.com/LandcLi/LandcOrg/api/invitation/v1"
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/service"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

func init() { di.Provide[invApi.InvitationController]("invitation.controller", &invCtrl{}) }

type invCtrl struct{}

func (c *invCtrl) CreateInvitation(ctx context.Context, req *v1.CreateInvitationRequest) (*v1.CreateInvitationResponse, error) {
	inv, err := service.GetInvitationService().CreateInvitation(
		req.InviterID, req.DepartmentID, req.PositionID, req.InviteeID, req.MaxUses, req.ExpiresInHours,
	)
	if err != nil {
		return nil, err
	}
	return &v1.CreateInvitationResponse{
		ID:        inv.ID,
		Code:      inv.Code,
		ExpiresAt: inv.ExpiresAt,
	}, nil
}

func (c *invCtrl) AcceptInvitation(ctx context.Context, req *v1.AcceptInvitationRequest) (*v1.AcceptInvitationResponse, error) {
	err := service.GetInvitationService().AcceptInvitation(req.UserID, req.Code)
	if err != nil {
		return nil, err
	}
	// Get department ID from invitation
	inv, _ := dao.GetInvitationDAO().GetByCode(req.Code)
	deptID := uint(0)
	if inv != nil {
		deptID = inv.DepartmentID
	}
	return &v1.AcceptInvitationResponse{Success: true, DepartmentID: deptID}, nil
}

func (c *invCtrl) RevokeInvitation(ctx context.Context, req *v1.RevokeInvitationRequest) (*v1.RevokeInvitationResponse, error) {
	if err := service.GetInvitationService().RevokeInvitation(req.InviterID, req.InvitationID); err != nil {
		return nil, err
	}
	return &v1.RevokeInvitationResponse{Success: true}, nil
}

func (c *invCtrl) ListInvitations(ctx context.Context, req *v1.ListInvitationsRequest) (*v1.ListInvitationsResponse, error) {
	invitations, total, err := service.GetInvitationService().ListMyInvitations(req.InviterID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]v1.InvitationInfo, 0, len(invitations))
	for _, inv := range invitations {
		info := v1.InvitationInfo{
			ID:           inv.ID,
			Code:         inv.Code,
			DepartmentID: inv.DepartmentID,
			PositionID:   inv.PositionID,
			InviteeID:    inv.InviteeID,
			Status:       inv.Status,
			MaxUses:      inv.MaxUses,
			UsedCount:    inv.UsedCount,
			ExpiresAt:    inv.ExpiresAt,
			CreatedAt:    inv.CreatedAt,
		}
		// Enrich department name
		dept, _ := dao.GetDepartmentDAO().GetByID(inv.DepartmentID)
		if dept != nil {
			info.DepartmentName = dept.Name
		}
		list = append(list, info)
	}

	return &v1.ListInvitationsResponse{Total: total, Invitations: list}, nil
}
