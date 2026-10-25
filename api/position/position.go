package position

import (
	"context"

	v1 "github.com/LandcLi/LandcOrg/api/position/v1"
	"github.com/LandcLi/landc-go/frame/pkg/di"
)

// PositionController defines the position management interface.
type PositionController interface {
	CreatePosition(ctx context.Context, req *v1.CreatePositionRequest) (*v1.CreatePositionResponse, error)
	UpdatePosition(ctx context.Context, req *v1.UpdatePositionRequest) (*v1.UpdatePositionResponse, error)
	DeletePosition(ctx context.Context, req *v1.DeletePositionRequest) (*v1.DeletePositionResponse, error)
	GetPosition(ctx context.Context, req *v1.GetPositionRequest) (*v1.GetPositionResponse, error)
	ListPositions(ctx context.Context, req *v1.ListPositionsRequest) (*v1.ListPositionsResponse, error)
}

var PositionGateway = di.NewGateway[PositionController]("position.controller")

func GetPositionController() PositionController {
	return PositionGateway.Get()
}
