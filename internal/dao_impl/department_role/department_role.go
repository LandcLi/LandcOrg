package department_role

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
	"github.com/go-redis/redis/v8"
)

func init() { dao.RegisterDepartmentRoleDAO(&deptRoleDAOImpl{redisClient: db.GetRedis()}) }

type deptRoleDAOImpl struct {
	redisClient *redis.Client
}

func (d *deptRoleDAOImpl) Assign(departmentID uint, roleIDs []uint) error {
	for _, roleID := range roleIDs {
		dr := &model.DepartmentRole{DepartmentID: departmentID, RoleID: roleID}
		if err := db.GetDB().Where("department_id = ? AND role_id = ?", departmentID, roleID).
			FirstOrCreate(dr).Error; err != nil {
			return err
		}
	}
	d.invalidateCache(departmentID)
	return nil
}

func (d *deptRoleDAOImpl) Revoke(departmentID uint, roleIDs []uint) error {
	err := db.GetDB().Where("department_id = ? AND role_id IN ?", departmentID, roleIDs).
		Delete(&model.DepartmentRole{}).Error
	if err == nil {
		d.invalidateCache(departmentID)
	}
	return err
}

func (d *deptRoleDAOImpl) GetRoleIDsByDepartmentID(departmentID uint) ([]uint, error) {
	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept_roles:%d", departmentID)
		data, err := d.redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var roleIDs []uint
			if err := json.Unmarshal(data, &roleIDs); err == nil {
				return roleIDs, nil
			}
		}
	}

	var roleIDs []uint
	err := db.GetDB().Model(&model.DepartmentRole{}).
		Where("department_id = ?", departmentID).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}

	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept_roles:%d", departmentID)
		data, _ := json.Marshal(roleIDs)
		d.redisClient.Set(context.Background(), cacheKey, data, 30*time.Minute)
	}

	return roleIDs, nil
}

func (d *deptRoleDAOImpl) GetRoleIDsByDepartmentIDs(departmentIDs []uint) ([]uint, error) {
	var roleIDs []uint
	if len(departmentIDs) == 0 {
		return roleIDs, nil
	}
	err := db.GetDB().Model(&model.DepartmentRole{}).
		Where("department_id IN ?", departmentIDs).
		Distinct("role_id").
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

func (d *deptRoleDAOImpl) invalidateCache(departmentID uint) {
	if d.redisClient == nil {
		return
	}
	d.redisClient.Del(context.Background(), fmt.Sprintf("dept_roles:%d", departmentID))
}
