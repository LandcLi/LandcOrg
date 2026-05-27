package service

import (
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

// PositionService defines the position business logic interface.
type PositionService interface {
	CreatePosition(name, code string, level int, description string) (*model.Position, error)
	UpdatePosition(id uint, name, code string, level *int, description string, status *int) error
	DeletePosition(id uint) error
	GetPosition(id uint) (*model.Position, error)
	ListPositions(page, pageSize int, keyword string) ([]model.Position, int64, error)
}

func RegisterPositionService(impl PositionService) {
	di.Provide[PositionService]("position.service", impl)
}

func GetPositionService() PositionService {
	return di.Require[PositionService]("position.service")
}
