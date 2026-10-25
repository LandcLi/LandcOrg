package sqlite_test

import (
	"encoding/json"
	"strconv"
	"testing"

	framedb "github.com/LandcLi/landc-go/frame/pkg/db"

	landcorg "github.com/LandcLi/LandcOrg"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

// initOrg 用临时 SQLite 完成一次 Org 初始化（含建表与 DI 注册）。
func initOrg(t *testing.T) {
	t.Helper()
	cfg := setupSQLiteEnv(t)
	if _, err := landcorg.InitWithConfig(cfg); err != nil {
		t.Fatalf("InitWithConfig: %v", err)
	}
}

func itoa(v uint) string { return strconv.FormatUint(uint64(v), 10) }

// newDept 创建部门并返回，失败直接 fail。
func newDept(t *testing.T, name, code string, parentID uint) *model.Department {
	t.Helper()
	d, err := service.GetDepartmentService().CreateDepartment(name, code, parentID, 0, "", "")
	if err != nil {
		t.Fatalf("create department %s: %v", code, err)
	}
	return d
}

func mustGetDept(t *testing.T, id uint) *model.Department {
	t.Helper()
	d, err := service.GetDepartmentService().GetDepartment(id)
	if err != nil {
		t.Fatalf("get department %d: %v", id, err)
	}
	return d
}

func uintPtr(v uint) *uint { return &v }

// TestMoveDepartmentRebasesSubtree 验证移动部门后整棵子树的 level 与 path 一起重算。
//
// 历史缺陷：UpdateDepartment 只改 ParentID，Level / Path 停留在旧值 ——
// 按 level 做数据权限、按 path 做前缀匹配都会算错，且没有任何报错。
func TestMoveDepartmentRebasesSubtree(t *testing.T) {
	initOrg(t)

	rootA := newDept(t, "研发中心", "rd", 0)
	childB := newDept(t, "后端组", "rd-be", rootA.ID)
	childC := newDept(t, "Go 小组", "rd-be-go", childB.ID)
	rootD := newDept(t, "产品中心", "pd", 0)

	if childB.Level != 2 || childC.Level != 3 {
		t.Fatalf("前置条件错误: B.level=%d C.level=%d, want 2 / 3", childB.Level, childC.Level)
	}

	// 把 B（连同其下的 C）整体移到 D 之下
	if err := service.GetDepartmentService().UpdateDepartment(childB.ID, "", "", uintPtr(rootD.ID), nil, "", "", nil); err != nil {
		t.Fatalf("move department: %v", err)
	}

	gotB := mustGetDept(t, childB.ID)
	wantBPath := "/" + itoa(rootD.ID) + "/"
	if gotB.ParentID != rootD.ID {
		t.Errorf("B.parent_id = %d, want %d", gotB.ParentID, rootD.ID)
	}
	if gotB.Level != 2 {
		t.Errorf("B.level = %d, want 2", gotB.Level)
	}
	if gotB.Path != wantBPath {
		t.Errorf("B.path = %q, want %q", gotB.Path, wantBPath)
	}

	gotC := mustGetDept(t, childC.ID)
	wantCPath := wantBPath + itoa(childB.ID) + "/"
	if gotC.Level != 3 {
		t.Errorf("子孙 C.level = %d, want 3（未随父级平移）", gotC.Level)
	}
	if gotC.Path != wantCPath {
		t.Errorf("子孙 C.path = %q, want %q", gotC.Path, wantCPath)
	}

	// 再移回根，层级应回到 1 / 2
	if err := service.GetDepartmentService().UpdateDepartment(childB.ID, "", "", uintPtr(0), nil, "", "", nil); err != nil {
		t.Fatalf("move to root: %v", err)
	}
	gotB = mustGetDept(t, childB.ID)
	gotC = mustGetDept(t, childC.ID)
	if gotB.Level != 1 || gotB.Path != "/" {
		t.Errorf("移回根后 B = {level:%d path:%q}, want {1, \"/\"}", gotB.Level, gotB.Path)
	}
	if gotC.Level != 2 || gotC.Path != "/"+itoa(childB.ID)+"/" {
		t.Errorf("移回根后 C = {level:%d path:%q}, want {2, \"/%s/\"}", gotC.Level, gotC.Path, itoa(childB.ID))
	}
}

// TestMoveDepartmentRejectsCycle 验证拒绝把部门挂到自己或自己的子孙下。
//
// 成环后 GetDepartmentTree 的父子互引会让 json 序列化因循环引用失败，
// 整棵树不可读 —— 必须挡在写入口。
func TestMoveDepartmentRejectsCycle(t *testing.T) {
	initOrg(t)

	root := newDept(t, "总部", "hq", 0)
	mid := newDept(t, "事业部", "bu", root.ID)
	leaf := newDept(t, "项目组", "team", mid.ID)

	cases := []struct {
		name   string
		target uint
	}{
		{"挂到自己", root.ID},
		{"挂到直接子部门", mid.ID},
		{"挂到孙部门", leaf.ID},
	}
	for _, tc := range cases {
		err := service.GetDepartmentService().UpdateDepartment(root.ID, "", "", uintPtr(tc.target), nil, "", "", nil)
		if err == nil {
			t.Errorf("%s: 应拒绝，实际成功", tc.name)
		}
	}

	// 被拒绝后数据必须保持原样
	if got := mustGetDept(t, root.ID); got.ParentID != 0 || got.Level != 1 {
		t.Errorf("拒绝后 root 被改动: parent=%d level=%d", got.ParentID, got.Level)
	}
}

// TestDeleteDepartmentRejectsMembers 验证有成员的部门不允许删除。
//
// 历史缺陷：只拦子部门不拦成员，删掉后 OrgMember.DepartmentID 变成悬空引用，
// 成员在「我的部门」与部门成员列表里都查不到，属于静默丢人。
func TestDeleteDepartmentRejectsMembers(t *testing.T) {
	initOrg(t)

	dept := newDept(t, "市场部", "mkt", 0)
	if err := service.GetMemberService().AddMember(1001, dept.ID, 0, false); err != nil {
		t.Fatalf("add member: %v", err)
	}

	if err := service.GetDepartmentService().DeleteDepartment(dept.ID); err == nil {
		t.Fatal("有成员的部门删除应被拒绝，实际成功")
	}

	// 移除成员后应可删除
	if err := service.GetMemberService().RemoveMember(1001, dept.ID); err != nil {
		t.Fatalf("remove member: %v", err)
	}
	if err := service.GetDepartmentService().DeleteDepartment(dept.ID); err != nil {
		t.Fatalf("清空成员后删除应成功: %v", err)
	}
}

// TestDeleteDepartmentCleansRoleBinding 验证删除部门时一并清掉角色绑定。
//
// 不清的话 getUserInheritedRoles 会继续把已删部门的角色算给用户 ——
// 部门都没了，权限还在。
func TestDeleteDepartmentCleansRoleBinding(t *testing.T) {
	initOrg(t)

	dept := newDept(t, "财务部", "fin", 0)
	if err := service.GetOrgRoleService().AssignDepartmentRoles(dept.ID, []uint{7, 8}); err != nil {
		t.Fatalf("assign roles: %v", err)
	}

	if err := service.GetDepartmentService().DeleteDepartment(dept.ID); err != nil {
		t.Fatalf("delete department: %v", err)
	}

	roleIDs, err := service.GetOrgRoleService().GetDepartmentRoles(dept.ID)
	if err != nil {
		t.Fatalf("get department roles after delete: %v", err)
	}
	if len(roleIDs) != 0 {
		t.Errorf("部门已删除，仍有角色绑定: %v", roleIDs)
	}
}

// TestTreeSurvivesCyclicData 验证历史脏数据（环）不会让部门树接口崩溃。
//
// 旧版本 UpdateDepartment 不校验父级，库里可能已存在 A.parent=B、B.parent=A。
// 不防御时 GetDepartmentTree 的父子互引会让 json.Marshal 失败，整棵树读不出来。
func TestTreeSurvivesCyclicData(t *testing.T) {
	initOrg(t)

	a := newDept(t, "环A", "cyc-a", 0)
	b := newDept(t, "环B", "cyc-b", a.ID)

	// 直接把父子互指写进库，绕过服务层校验伪造环
	db := framedb.GetDB()
	if err := db.Model(&model.Department{}).Where("id = ?", a.ID).Update("parent_id", b.ID).Error; err != nil {
		t.Fatalf("corrupt A: %v", err)
	}
	if err := db.Model(&model.Department{}).Where("id = ?", b.ID).Update("parent_id", a.ID).Error; err != nil {
		t.Fatalf("corrupt B: %v", err)
	}

	tree, err := service.GetDepartmentService().GetDepartmentTree(0)
	if err != nil {
		t.Fatalf("GetDepartmentTree on cyclic data: %v", err)
	}
	if _, err := json.Marshal(tree); err != nil {
		t.Fatalf("环数据导致树无法序列化: %v", err)
	}
	if len(tree) == 0 {
		t.Error("环上节点应作为根返回，不能凭空消失")
	}
}
