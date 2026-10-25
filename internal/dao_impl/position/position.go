package position

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterPositionDAO(&posDAOImpl{}) }

type posDAOImpl struct{}

func (d *posDAOImpl) Create(pos *model.Position) error { return db.GetDB().Create(pos).Error }
func (d *posDAOImpl) Update(pos *model.Position) error { return db.GetDB().Save(pos).Error }
func (d *posDAOImpl) Delete(id uint) error             { return db.GetDB().Delete(&model.Position{}, id).Error }

func (d *posDAOImpl) GetByID(id uint) (*model.Position, error) {
	var pos model.Position
	err := db.GetDB().First(&pos, id).Error
	return &pos, err
}

// GetByIDs 批量按主键查询，一次 IN 查询替代逐个 GetByID。
func (d *posDAOImpl) GetByIDs(ids []uint) ([]model.Position, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var positions []model.Position
	err := db.GetDB().Where("id IN ?", ids).Find(&positions).Error
	return positions, err
}

func (d *posDAOImpl) GetByCode(code string) (*model.Position, error) {
	var pos model.Position
	err := db.GetDB().Where("code = ?", code).First(&pos).Error
	return &pos, err
}

func (d *posDAOImpl) List(page, pageSize int, keyword string) ([]model.Position, int64, error) {
	var positions []model.Position
	var total int64
	query := db.GetDB().Model(&model.Position{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err := query.Find(&positions).Error
	return positions, total, err
}
