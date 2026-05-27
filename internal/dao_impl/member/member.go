package member

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterMemberDAO(&memberDAOImpl{}) }

type memberDAOImpl struct{}

func (d *memberDAOImpl) Create(m *model.OrgMember) error { return db.GetDB().Create(m).Error }
func (d *memberDAOImpl) Update(m *model.OrgMember) error { return db.GetDB().Save(m).Error }

func (d *memberDAOImpl) Delete(userID, departmentID uint) error {
	return db.GetDB().Where("user_id = ? AND department_id = ?", userID, departmentID).
		Delete(&model.OrgMember{}).Error
}

func (d *memberDAOImpl) GetByUserAndDept(userID, departmentID uint) (*model.OrgMember, error) {
	var m model.OrgMember
	err := db.GetDB().Where("user_id = ? AND department_id = ?", userID, departmentID).First(&m).Error
	return &m, err
}

func (d *memberDAOImpl) GetByUserID(userID uint) ([]model.OrgMember, error) {
	var members []model.OrgMember
	err := db.GetDB().Where("user_id = ?", userID).Find(&members).Error
	return members, err
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
