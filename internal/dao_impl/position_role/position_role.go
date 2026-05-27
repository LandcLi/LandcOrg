package position_role

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterPositionRoleDAO(&posRoleDAOImpl{}) }

type posRoleDAOImpl struct{}

func (d *posRoleDAOImpl) Assign(positionID uint, roleIDs []uint) error {
	for _, roleID := range roleIDs {
		pr := &model.PositionRole{PositionID: positionID, RoleID: roleID}
		if err := db.GetDB().Where("position_id = ? AND role_id = ?", positionID, roleID).
			FirstOrCreate(pr).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *posRoleDAOImpl) Revoke(positionID uint, roleIDs []uint) error {
	return db.GetDB().Where("position_id = ? AND role_id IN ?", positionID, roleIDs).
		Delete(&model.PositionRole{}).Error
}

func (d *posRoleDAOImpl) GetRoleIDsByPositionID(positionID uint) ([]uint, error) {
	var roleIDs []uint
	err := db.GetDB().Model(&model.PositionRole{}).
		Where("position_id = ?", positionID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (d *posRoleDAOImpl) GetRoleIDsByPositionIDs(positionIDs []uint) ([]uint, error) {
	var roleIDs []uint
	if len(positionIDs) == 0 {
		return roleIDs, nil
	}
	err := db.GetDB().Model(&model.PositionRole{}).
		Where("position_id IN ?", positionIDs).
		Distinct("role_id").
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}
