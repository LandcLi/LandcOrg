package position

import (
	"context"

	posApi "github.com/LandcLi/LandcOrg/api/position"
	v1 "github.com/LandcLi/LandcOrg/api/position/v1"
	"github.com/LandcLi/LandcOrg/service"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

func init() { di.Provide[posApi.PositionController]("position.controller", &posCtrl{}) }

type posCtrl struct{}

func (c *posCtrl) CreatePosition(ctx context.Context, req *v1.CreatePositionRequest) (*v1.CreatePositionResponse, error) {
	pos, err := service.GetPositionService().CreatePosition(req.Name, req.Code, req.Level, req.Description)
	if err != nil { return nil, err }
	return &v1.CreatePositionResponse{ID: pos.ID, Name: pos.Name, Code: pos.Code}, nil
}

func (c *posCtrl) UpdatePosition(ctx context.Context, req *v1.UpdatePositionRequest) (*v1.UpdatePositionResponse, error) {
	err := service.GetPositionService().UpdatePosition(req.ID, req.Name, req.Code, req.Level, req.Description, req.Status)
	if err != nil { return nil, err }
	return &v1.UpdatePositionResponse{Success: true}, nil
}

func (c *posCtrl) DeletePosition(ctx context.Context, req *v1.DeletePositionRequest) (*v1.DeletePositionResponse, error) {
	if err := service.GetPositionService().DeletePosition(req.ID); err != nil { return nil, err }
	return &v1.DeletePositionResponse{Success: true}, nil
}

func (c *posCtrl) GetPosition(ctx context.Context, req *v1.GetPositionRequest) (*v1.GetPositionResponse, error) {
	pos, err := service.GetPositionService().GetPosition(req.ID)
	if err != nil { return nil, err }
	return &v1.GetPositionResponse{ID: pos.ID, Name: pos.Name, Code: pos.Code, Level: pos.Level, Description: pos.Description, Status: pos.Status}, nil
}

func (c *posCtrl) ListPositions(ctx context.Context, req *v1.ListPositionsRequest) (*v1.ListPositionsResponse, error) {
	positions, total, err := service.GetPositionService().ListPositions(req.Page, req.PageSize, req.Keyword)
	if err != nil { return nil, err }
	list := make([]v1.GetPositionResponse, 0, len(positions))
	for _, p := range positions {
		list = append(list, v1.GetPositionResponse{ID: p.ID, Name: p.Name, Code: p.Code, Level: p.Level, Description: p.Description, Status: p.Status})
	}
	return &v1.ListPositionsResponse{Total: total, List: list}, nil
}
