package position_role

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterPositionRoleDAO(&posRoleDAOImpl{}) }

type posRoleDAOImpl struct {
}

func (d *posRoleDAOImpl) Assign(positionID uint, roleIDs []uint) error {
	for _, roleID := range roleIDs {
		pr := &model.PositionRole{PositionID: positionID, RoleID: roleID}
		if err := db.GetDB().Where("position_id = ? AND role_id = ?", positionID, roleID).
			FirstOrCreate(pr).Error; err != nil {
			return err
		}
	}
	d.invalidateCache(positionID)
	return nil
}

func (d *posRoleDAOImpl) Revoke(positionID uint, roleIDs []uint) error {
	err := db.GetDB().Where("position_id = ? AND role_id IN ?", positionID, roleIDs).
		Delete(&model.PositionRole{}).Error
	if err == nil {
		d.invalidateCache(positionID)
	}
	return err
}

func (d *posRoleDAOImpl) GetRoleIDsByPositionID(positionID uint) ([]uint, error) {
	if db.GetRedis() != nil {
		cacheKey := fmt.Sprintf("pos_roles:%d", positionID)
		data, err := db.GetRedis().Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var roleIDs []uint
			if err := json.Unmarshal(data, &roleIDs); err == nil {
				return roleIDs, nil
			}
		}
	}

	var roleIDs []uint
	err := db.GetDB().Model(&model.PositionRole{}).
		Where("position_id = ?", positionID).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}

	if db.GetRedis() != nil {
		cacheKey := fmt.Sprintf("pos_roles:%d", positionID)
		data, _ := json.Marshal(roleIDs)
		db.GetRedis().Set(context.Background(), cacheKey, data, 30*time.Minute)
	}

	return roleIDs, nil
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

func (d *posRoleDAOImpl) invalidateCache(positionID uint) {
	if db.GetRedis() == nil {
		return
	}
	db.GetRedis().Del(context.Background(), fmt.Sprintf("pos_roles:%d", positionID))
}
