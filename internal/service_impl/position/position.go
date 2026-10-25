package position

import (
	"fmt"

	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterPositionService(&posServiceImpl{}) }

type posServiceImpl struct{}

func (s *posServiceImpl) CreatePosition(name, code string, level int, description string) (*model.Position, error) {
	existing, _ := dao.GetPositionDAO().GetByCode(code)
	if existing != nil && existing.ID != 0 {
		return nil, fmt.Errorf("position code '%s' already exists", code)
	}
	pos := &model.Position{Name: name, Code: code, Level: level, Description: description, Status: 1}
	if err := dao.GetPositionDAO().Create(pos); err != nil {
		return nil, err
	}
	return pos, nil
}

func (s *posServiceImpl) UpdatePosition(id uint, name, code string, level *int, description string, status *int) error {
	pos, err := dao.GetPositionDAO().GetByID(id)
	if err != nil {
		return fmt.Errorf("position not found: %w", err)
	}
	if name != "" {
		pos.Name = name
	}
	if code != "" {
		existing, _ := dao.GetPositionDAO().GetByCode(code)
		if existing != nil && existing.ID != id {
			return fmt.Errorf("position code '%s' already exists", code)
		}
		pos.Code = code
	}
	if level != nil {
		pos.Level = *level
	}
	if description != "" {
		pos.Description = description
	}
	if status != nil {
		pos.Status = *status
	}
	return dao.GetPositionDAO().Update(pos)
}

func (s *posServiceImpl) DeletePosition(id uint) error { return dao.GetPositionDAO().Delete(id) }
func (s *posServiceImpl) GetPosition(id uint) (*model.Position, error) {
	return dao.GetPositionDAO().GetByID(id)
}

func (s *posServiceImpl) ListPositions(page, pageSize int, keyword string) ([]model.Position, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return dao.GetPositionDAO().List(page, pageSize, keyword)
}
