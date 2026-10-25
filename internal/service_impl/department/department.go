package department

import (
	"fmt"
	"strconv"
	"strings"

	v1 "github.com/LandcLi/LandcOrg/api/department/v1"
	"github.com/LandcLi/LandcOrg/dao"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

func init() { service.RegisterDepartmentService(&deptServiceImpl{}) }

type deptServiceImpl struct{}

func (s *deptServiceImpl) CreateDepartment(name, code string, parentID uint, sort int, leader, description string) (*model.Department, error) {
	existing, _ := dao.GetDepartmentDAO().GetByCode(code)
	if existing != nil && existing.ID != 0 {
		return nil, fmt.Errorf("department code '%s' already exists", code)
	}

	dept := &model.Department{
		Name:        name,
		Code:        code,
		ParentID:    parentID,
		Sort:        sort,
		Leader:      leader,
		Description: description,
		Status:      1,
		Level:       1,
	}

	// Calculate level and path
	if parentID > 0 {
		parent, err := dao.GetDepartmentDAO().GetByID(parentID)
		if err != nil {
			return nil, fmt.Errorf("parent department not found: %w", err)
		}
		dept.Level = parent.Level + 1
		dept.Path = childPath(parent)
	} else {
		dept.Path = "/"
	}

	if err := dao.GetDepartmentDAO().Create(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// childPath 返回「挂在 dept 之下」的部门应有的 path 前缀。
// 根部门 path 为 "/"，其子部门 path 为 "/<根ID>/"，因此统一取 dept.Path + dept.ID + "/"。
func childPath(dept *model.Department) string {
	return dept.Path + strconv.FormatUint(uint64(dept.ID), 10) + "/"
}

func (s *deptServiceImpl) UpdateDepartment(id uint, name, code string, parentID *uint, sort *int, leader, description string, status *int) error {
	dept, err := dao.GetDepartmentDAO().GetByID(id)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}
	if name != "" {
		dept.Name = name
	}
	if code != "" {
		existing, _ := dao.GetDepartmentDAO().GetByCode(code)
		if existing != nil && existing.ID != id {
			return fmt.Errorf("department code '%s' already exists", code)
		}
		dept.Code = code
	}
	if sort != nil {
		dept.Sort = *sort
	}
	if leader != "" {
		dept.Leader = leader
	}
	if description != "" {
		dept.Description = description
	}
	if status != nil {
		dept.Status = *status
	}

	// 移动部门：先校验目标父级不会造成环，再把 level/path 连同整个子树一起重算。
	//
	// 历史缺陷：此前只是把 ParentID 改掉就落库，Level / Path 停留在旧值。
	// 后果是子孙部门的层级与路径永久失真 —— 按 level 做数据权限、按 path
	// 做前缀匹配的功能会全部算错，而且没有任何报错线索。
	if parentID != nil && *parentID != dept.ParentID {
		if err := s.validateMove(id, *parentID); err != nil {
			return err
		}

		oldLevel := dept.Level
		oldChildPrefix := childPath(dept)

		if *parentID == 0 {
			dept.ParentID = 0
			dept.Level = 1
			dept.Path = "/"
		} else {
			parent, err := dao.GetDepartmentDAO().GetByID(*parentID)
			if err != nil {
				return fmt.Errorf("parent department not found: %w", err)
			}
			dept.ParentID = *parentID
			dept.Level = parent.Level + 1
			dept.Path = childPath(parent)
		}

		if err := dao.GetDepartmentDAO().Update(dept); err != nil {
			return err
		}
		return s.rebaseSubtree(oldChildPrefix, childPath(dept), dept.Level-oldLevel)
	}

	return dao.GetDepartmentDAO().Update(dept)
}

// validateMove 校验把 deptID 挂到 newParentID 之下是否合法。
//
// 不仅检查「不能挂到自己」，还要沿目标父级的祖先链上溯，防止挂到自己的
// 子孙下 —— 那样会形成 A→B→A 的环。成环后 GetDepartmentTree 序列化时
// 父子互相引用，json 会因循环引用直接失败，整棵树不可读。
func (s *deptServiceImpl) validateMove(deptID, newParentID uint) error {
	if newParentID == 0 {
		return nil
	}
	if newParentID == deptID {
		return fmt.Errorf("cannot move department %d under itself", deptID)
	}

	seen := make(map[uint]struct{})
	for cur := newParentID; cur != 0; {
		if cur == deptID {
			return fmt.Errorf("cannot move department %d under its own descendant %d", deptID, newParentID)
		}
		if _, dup := seen[cur]; dup {
			// 目标父级自身已处于环中，任何挂载都无意义，直接拒绝并暴露脏数据。
			return fmt.Errorf("department %d is already in a cycle, fix the data before moving", newParentID)
		}
		seen[cur] = struct{}{}

		parent, err := dao.GetDepartmentDAO().GetByID(cur)
		if err != nil {
			return fmt.Errorf("parent department %d not found: %w", cur, err)
		}
		cur = parent.ParentID
	}
	return nil
}

// rebaseSubtree 把 oldPrefix 下的所有子孙路径改写为 newPrefix，并整体平移层级。
//
// dept 自身的 level 变化量以 delta 传入，子孙的 level 需要同步平移同样的量，
// 否则移动后子部门会比父部门「深」或「浅」得莫名其妙。
func (s *deptServiceImpl) rebaseSubtree(oldPrefix, newPrefix string, levelDelta int) error {
	if oldPrefix == newPrefix {
		return nil
	}

	all, err := dao.GetDepartmentDAO().GetAll()
	if err != nil {
		return err
	}

	for i := range all {
		d := &all[i]
		if !strings.HasPrefix(d.Path, oldPrefix) {
			continue
		}
		d.Path = newPrefix + strings.TrimPrefix(d.Path, oldPrefix)
		// 整棵子树随父级一起升降，否则移动后子部门可能比父部门还「浅」。
		d.Level += levelDelta
		if err := dao.GetDepartmentDAO().Update(d); err != nil {
			return fmt.Errorf("rebase department %d: %w", d.ID, err)
		}
	}
	return nil
}

func (s *deptServiceImpl) DeleteDepartment(id uint) error {
	// Check if has children
	children, _ := dao.GetDepartmentDAO().GetByParentID(id)
	if len(children) > 0 {
		return fmt.Errorf("cannot delete department with sub-departments")
	}

	// 有成员的部门不允许直接删除。
	// 历史缺陷：此前只拦子部门，删掉后 OrgMember.DepartmentID 变成悬空引用，
	// 成员在「我的部门」里查不到、在部门成员列表里也查不到，属于静默丢人。
	if _, total, err := dao.GetMemberDAO().GetByDepartmentID(id, 0, 0); err == nil && total > 0 {
		return fmt.Errorf("cannot delete department with %d member(s), remove or transfer them first", total)
	}

	if err := dao.GetDepartmentDAO().Delete(id); err != nil {
		return err
	}

	// 部门没了，绑定在它上面的角色必须一起清掉，否则 getUserInheritedRoles
	// 会继续把这些角色算给用户 —— 部门都删了权限还在。
	if err := dao.GetDepartmentRoleDAO().DeleteByDepartmentID(id); err != nil {
		return fmt.Errorf("cleanup department roles: %w", err)
	}
	return nil
}

func (s *deptServiceImpl) GetDepartment(id uint) (*model.Department, error) {
	return dao.GetDepartmentDAO().GetByID(id)
}

func (s *deptServiceImpl) GetDepartmentTree(rootID uint) ([]*v1.DepartmentTreeNode, error) {
	allDepts, err := dao.GetDepartmentDAO().GetAll()
	if err != nil {
		return nil, err
	}

	// Build tree
	nodeMap := make(map[uint]*v1.DepartmentTreeNode)
	for _, d := range allDepts {
		nodeMap[d.ID] = &v1.DepartmentTreeNode{
			ID:       d.ID,
			Name:     d.Name,
			Code:     d.Code,
			ParentID: d.ParentID,
			Level:    d.Level,
			Sort:     d.Sort,
			Leader:   d.Leader,
			Status:   d.Status,
			Children: make([]*v1.DepartmentTreeNode, 0),
		}
	}

	// 历史脏数据里可能存在环（旧版本 UpdateDepartment 不校验父级）。
	// 一旦成环，父子节点会互相引用，json 序列化会因循环引用失败，整棵树不可读。
	// 这里把环上的节点当作根节点处理，保证接口永远能返回一棵合法的树。
	byID := make(map[uint]*model.Department, len(allDepts))
	for i := range allDepts {
		byID[allDepts[i].ID] = &allDepts[i]
	}
	cyclic := detectCycles(byID)

	var roots []*v1.DepartmentTreeNode
	for _, d := range allDepts {
		node := nodeMap[d.ID]
		switch {
		case cyclic[d.ID]:
			// 环上节点不挂父级，直接作为根，避免无限递归。
			roots = append(roots, node)
		case d.ParentID == rootID:
			roots = append(roots, node)
		default:
			parent, ok := nodeMap[d.ParentID]
			if !ok || cyclic[d.ParentID] {
				// 父级不存在（孤儿）或父级在环上，挂回根节点，不让它凭空消失。
				roots = append(roots, node)
				continue
			}
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// detectCycles 找出所有处在环上的部门。
// 判断方式：从节点沿 ParentID 上溯，若在到达根之前回到已访问过的节点，则该链上的节点成环。
func detectCycles(byID map[uint]*model.Department) map[uint]bool {
	cyclic := make(map[uint]bool)
	for id := range byID {
		path := make(map[uint]struct{})
		for cur := id; cur != 0; {
			if _, ok := path[cur]; ok {
				// 从 cur 开始往回都是环成员
				for m := range path {
					cyclic[m] = true
				}
				break
			}
			path[cur] = struct{}{}
			d, ok := byID[cur]
			if !ok {
				break // 父级缺失，不是环而是孤儿
			}
			cur = d.ParentID
		}
	}
	return cyclic
}

func (s *deptServiceImpl) ListDepartments(page, pageSize int, parentID *uint, keyword string) ([]model.Department, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return dao.GetDepartmentDAO().List(page, pageSize, parentID, keyword)
}
