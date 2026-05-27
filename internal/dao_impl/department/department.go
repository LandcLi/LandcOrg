package department

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

func init() { dao.RegisterDepartmentDAO(&deptDAOImpl{redisClient: db.GetRedis()}) }

type deptDAOImpl struct {
	redisClient *redis.Client
}

func (d *deptDAOImpl) Create(dept *model.Department) error { return db.GetDB().Create(dept).Error }

func (d *deptDAOImpl) Update(dept *model.Department) error {
	err := db.GetDB().Save(dept).Error
	if err == nil {
		d.invalidateCache(dept.ID, dept.Code)
	}
	return err
}

func (d *deptDAOImpl) Delete(id uint) error {
	if d.redisClient != nil {
		if dept, err := d.GetByID(id); err == nil {
			defer d.invalidateCache(dept.ID, dept.Code)
		}
	}
	return db.GetDB().Delete(&model.Department{}, id).Error
}

func (d *deptDAOImpl) GetByID(id uint) (*model.Department, error) {
	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept:id:%d", id)
		data, err := d.redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var dept model.Department
			if err := json.Unmarshal(data, &dept); err == nil {
				return &dept, nil
			}
		}
	}

	var dept model.Department
	err := db.GetDB().First(&dept, id).Error
	if err != nil {
		return nil, err
	}

	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept:id:%d", id)
		data, _ := json.Marshal(dept)
		d.redisClient.Set(context.Background(), cacheKey, data, time.Hour)
	}

	return &dept, nil
}

func (d *deptDAOImpl) GetByCode(code string) (*model.Department, error) {
	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept:code:%s", code)
		data, err := d.redisClient.Get(context.Background(), cacheKey).Bytes()
		if err == nil {
			var dept model.Department
			if err := json.Unmarshal(data, &dept); err == nil {
				return &dept, nil
			}
		}
	}

	var dept model.Department
	err := db.GetDB().Where("code = ?", code).First(&dept).Error
	if err != nil {
		return nil, err
	}

	if d.redisClient != nil {
		cacheKey := fmt.Sprintf("dept:code:%s", code)
		data, _ := json.Marshal(dept)
		d.redisClient.Set(context.Background(), cacheKey, data, time.Hour)
	}

	return &dept, nil
}

func (d *deptDAOImpl) GetByParentID(parentID uint) ([]model.Department, error) {
	var depts []model.Department
	err := db.GetDB().Where("parent_id = ?", parentID).Order("sort ASC").Find(&depts).Error
	return depts, err
}

func (d *deptDAOImpl) GetAll() ([]model.Department, error) {
	var depts []model.Department
	err := db.GetDB().Order("sort ASC").Find(&depts).Error
	return depts, err
}

func (d *deptDAOImpl) List(page, pageSize int, parentID *uint, keyword string) ([]model.Department, int64, error) {
	var depts []model.Department
	var total int64
	query := db.GetDB().Model(&model.Department{})
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err := query.Order("sort ASC").Find(&depts).Error
	return depts, total, err
}

func (d *deptDAOImpl) invalidateCache(id uint, code string) {
	if d.redisClient == nil {
		return
	}
	ctx := context.Background()
	d.redisClient.Del(ctx, fmt.Sprintf("dept:id:%d", id))
	if code != "" {
		d.redisClient.Del(ctx, fmt.Sprintf("dept:code:%s", code))
	}
}
