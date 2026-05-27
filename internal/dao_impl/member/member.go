package member

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

func init() { dao.RegisterMemberDAO(&memberDAOImpl{redisClient: db.GetRedis()}) }

type memberDAOImpl struct {
	redisClient *redis.Client
}

func (d *memberDAOImpl) Create(m *model.OrgMember) error {
	err := db.GetDB().Create(m).Error
	if err == nil {
		d.invalidateUserCache(m.UserID)
	}
	return err
}

func (d *memberDAOImpl) Update(m *model.OrgMember) error {
	err := db.GetDB().Save(m).Error
	if err == nil {
		d.invalidateUserCache(m.UserID)
	}
	return err
}

func (d *memberDAOImpl) Delete(userID, departmentID uint) error {
	err := db.GetDB().Where("user_id = ? AND department_id = ?", userID, departmentID).
		Delete(&model.OrgMember{}).Error
	if err == nil {
		d.invalidateUserCache(userID)
	}
	return err
}

func (d *memberDAOImpl) GetByUserAndDept(userID, departmentID uint) (*model.OrgMember, error) {
	var m model.OrgMember
	err := db.GetDB().Where("user_id = ? AND department_id = ?", userID, departmentID).First(&m).Error
	return &m, err
}

func (d *memberDAOImpl) GetByUserID(userID uint) ([]model.OrgMember, error) {
	// Try cache
	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("member:user:%d", userID)
		data, err := d.redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var members []model.OrgMember
			if err := json.Unmarshal(data, &members); err == nil {
				return members, nil
			}
		}
	}

	var members []model.OrgMember
	err := db.GetDB().Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return nil, err
	}

	// Store in cache
	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("member:user:%d", userID)
		data, _ := json.Marshal(members)
		d.redisClient.Set(context.Background(), cacheKey, data, 30*time.Minute)
	}

	return members, nil
}

func (d *memberDAOImpl) GetByDepartmentID(departmentID uint, page, pageSize int) ([]model.OrgMember, int64, error) {
	var members []model.OrgMember
	var total int64
	query := db.GetDB().Model(&model.OrgMember{}).Where("department_id = ?", departmentID)
	query.Count(&total)
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err := query.Find(&members).Error
	return members, total, err
}

func (d *memberDAOImpl) invalidateUserCache(userID uint) {
	if d.redisClient == nil {
		return
	}
	d.redisClient.Del(context.Background(), fmt.Sprintf("member:user:%d", userID))
}
