package department_role

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterDepartmentRoleDAO(&deptRoleDAOImpl{}) }

type deptRoleDAOImpl struct{}

func (d *deptRoleDAOImpl) Assign(departmentID uint, roleIDs []uint) error {
	for _, roleID := range roleIDs {
		dr := &model.DepartmentRole{DepartmentID: departmentID, RoleID: roleID}
		// Ignore duplicate
		if err := db.GetDB().Where("department_id = ? AND role_id = ?", departmentID, roleID).
			FirstOrCreate(dr).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *deptRoleDAOImpl) Revoke(departmentID uint, roleIDs []uint) error {
	return db.GetDB().Where("department_id = ? AND role_id IN ?", departmentID, roleIDs).
		Delete(&model.DepartmentRole{}).Error
}

func (d *deptRoleDAOImpl) GetRoleIDsByDepartmentID(departmentID uint) ([]uint, error) {
	var roleIDs []uint
	err := db.GetDB().Model(&model.DepartmentRole{}).
		Where("department_id = ?", departmentID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
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
