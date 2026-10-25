package dao

import (
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/landc-go/frame/pkg/di"
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
func GetDepartmentDAO() DepartmentDAO          { return di.Require[DepartmentDAO]("department.dao") }

type PositionDAO interface {
	Create(pos *model.Position) error
	Update(pos *model.Position) error
	Delete(id uint) error
	GetByID(id uint) (*model.Position, error)
	GetByCode(code string) (*model.Position, error)
	List(page, pageSize int, keyword string) ([]model.Position, int64, error)
}

func RegisterPositionDAO(impl PositionDAO) { di.Provide[PositionDAO]("position.dao", impl) }
func GetPositionDAO() PositionDAO          { return di.Require[PositionDAO]("position.dao") }

type MemberDAO interface {
	Create(member *model.OrgMember) error
	Update(member *model.OrgMember) error
	Delete(userID, departmentID uint) error
	GetByUserAndDept(userID, departmentID uint) (*model.OrgMember, error)
	GetByUserID(userID uint) ([]model.OrgMember, error)
	GetByDepartmentID(departmentID uint, page, pageSize int) ([]model.OrgMember, int64, error)
}

func RegisterMemberDAO(impl MemberDAO) { di.Provide[MemberDAO]("member.dao", impl) }
func GetMemberDAO() MemberDAO          { return di.Require[MemberDAO]("member.dao") }

type InvitationDAO interface {
	Create(inv *model.OrgInvitation) error
	Update(inv *model.OrgInvitation) error
	GetByCode(code string) (*model.OrgInvitation, error)
	GetByID(id uint) (*model.OrgInvitation, error)
	ListByInviter(inviterID uint, page, pageSize int) ([]model.OrgInvitation, int64, error)
	ListByDepartment(departmentID uint, page, pageSize int) ([]model.OrgInvitation, int64, error)
}

func RegisterInvitationDAO(impl InvitationDAO) { di.Provide[InvitationDAO]("invitation.dao", impl) }
func GetInvitationDAO() InvitationDAO          { return di.Require[InvitationDAO]("invitation.dao") }

type DepartmentRoleDAO interface {
	Assign(departmentID uint, roleIDs []uint) error
	Revoke(departmentID uint, roleIDs []uint) error
	GetRoleIDsByDepartmentID(departmentID uint) ([]uint, error)
	GetRoleIDsByDepartmentIDs(departmentIDs []uint) ([]uint, error)
	// DeleteByDepartmentID 删除某部门的全部角色绑定。
	// 部门被删除时必须调用，否则残留绑定会让 getUserInheritedRoles 把已删部门的
	// 角色继续算给用户，形成「部门没了、权限还在」的幽灵授权。
	DeleteByDepartmentID(departmentID uint) error
}

func RegisterDepartmentRoleDAO(impl DepartmentRoleDAO) {
	di.Provide[DepartmentRoleDAO]("department_role.dao", impl)
}
func GetDepartmentRoleDAO() DepartmentRoleDAO {
	return di.Require[DepartmentRoleDAO]("department_role.dao")
}

type PositionRoleDAO interface {
	Assign(positionID uint, roleIDs []uint) error
	Revoke(positionID uint, roleIDs []uint) error
	GetRoleIDsByPositionID(positionID uint) ([]uint, error)
	GetRoleIDsByPositionIDs(positionIDs []uint) ([]uint, error)
}

func RegisterPositionRoleDAO(impl PositionRoleDAO) {
	di.Provide[PositionRoleDAO]("position_role.dao", impl)
}
func GetPositionRoleDAO() PositionRoleDAO {
	return di.Require[PositionRoleDAO]("position_role.dao")
}
