package invitation

import (
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/db"
)

func init() { dao.RegisterInvitationDAO(&invitationDAOImpl{}) }

type invitationDAOImpl struct{}

func (d *invitationDAOImpl) Create(inv *model.OrgInvitation) error {
	return db.GetDB().Create(inv).Error
}

func (d *invitationDAOImpl) Update(inv *model.OrgInvitation) error {
	return db.GetDB().Save(inv).Error
}

func (d *invitationDAOImpl) GetByCode(code string) (*model.OrgInvitation, error) {
	var inv model.OrgInvitation
	err := db.GetDB().Where("code = ?", code).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (d *invitationDAOImpl) GetByID(id uint) (*model.OrgInvitation, error) {
	var inv model.OrgInvitation
	err := db.GetDB().First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (d *invitationDAOImpl) ListByInviter(inviterID uint, page, pageSize int) ([]model.OrgInvitation, int64, error) {
	var list []model.OrgInvitation
	var total int64
	query := db.GetDB().Model(&model.OrgInvitation{}).Where("inviter_id = ?", inviterID)
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (d *invitationDAOImpl) ListByDepartment(departmentID uint, page, pageSize int) ([]model.OrgInvitation, int64, error) {
	var list []model.OrgInvitation
	var total int64
	query := db.GetDB().Model(&model.OrgInvitation{}).Where("department_id = ?", departmentID)
	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
