package dao

import (
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/mvc/pkg/di"
)

type DepartmentDAO interface {
	Create(dept *model.Department) error
	Update(dept *model.Department) error
	Delete(id uint) error
	GetByID(id uint) (*model.Department, error)
	GetByCode(code string) (*model.Department, error)
	GetByParentID(parentID uint) ([]model.Department, error)
	GetAll() ([]model.Department, error)
	List(page, pageSize int, parentID *uint, keyword string) ([]model.Department, int64, error)
}

func RegisterDepartmentDAO(impl DepartmentDAO) { di.Provide[DepartmentDAO]("department.dao", impl) }
func GetDepartmentDAO() DepartmentDAO           { return di.Require[DepartmentDAO]("department.dao") }

type PositionDAO interface {
	Create(pos *model.Position) error
	Update(pos *model.Position) error
	Delete(id uint) error
	GetByID(id uint) (*model.Position, error)
	GetByCode(code string) (*model.Position, error)
	List(page, pageSize int, keyword string) ([]model.Position, int64, error)
}

func RegisterPositionDAO(impl PositionDAO) { di.Provide[PositionDAO]("position.dao", impl) }
func GetPositionDAO() PositionDAO           { return di.Require[PositionDAO]("position.dao") }

type MemberDAO interface {
	Create(member *model.OrgMember) error
	Update(member *model.OrgMember) error
	Delete(userID, departmentID uint) error
	GetByUserAndDept(userID, departmentID uint) (*model.OrgMember, error)
	GetByUserID(userID uint) ([]model.OrgMember, error)
	GetByDepartmentID(departmentID uint, page, pageSize int) ([]model.OrgMember, int64, error)
}

func RegisterMemberDAO(impl MemberDAO) { di.Provide[MemberDAO]("member.dao", impl) }
func GetMemberDAO() MemberDAO           { return di.Require[MemberDAO]("member.dao") }
