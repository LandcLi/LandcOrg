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

var PositionRoutes = di.Routes{
	"CreatePosition": {Method: "POST", Path: "/api/position/create"},
	"UpdatePosition": {Method: "POST", Path: "/api/position/update"},
	"DeletePosition": {Method: "POST", Path: "/api/position/delete"},
	"GetPosition":    {Method: "POST", Path: "/api/position/get"},
	"ListPositions":  {Method: "POST", Path: "/api/position/list"},
}

var PositionGateway = di.NewGateway[PositionController]("position.controller")

func GetPositionController() PositionController {
	return PositionGateway.Get()
}
