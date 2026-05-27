package department

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterDepartmentDAO(&deptDAOImpl{}) }

type deptDAOImpl struct{}

func (d *deptDAOImpl) Create(dept *model.Department) error    { return db.GetDB().Create(dept).Error }
func (d *deptDAOImpl) Update(dept *model.Department) error    { return db.GetDB().Save(dept).Error }
func (d *deptDAOImpl) Delete(id uint) error                   { return db.GetDB().Delete(&model.Department{}, id).Error }

func (d *deptDAOImpl) GetByID(id uint) (*model.Department, error) {
	var dept model.Department
	err := db.GetDB().First(&dept, id).Error
	return &dept, err
}

func (d *deptDAOImpl) GetByCode(code string) (*model.Department, error) {
	var dept model.Department
	err := db.GetDB().Where("code = ?", code).First(&dept).Error
	return &dept, err
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
