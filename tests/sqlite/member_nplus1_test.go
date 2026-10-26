package sqlite_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm/logger"

	framedb "github.com/LandcLi/landc-go/frame/pkg/db"

	landcorg "github.com/LandcLi/LandcOrg"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/service"
)

// countingLogger 统计执行过的 SELECT 语句次数。
type countingLogger struct {
	logger.Interface
	selects int
}

func (c *countingLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 直接调用底层 logger 打印（保持原有日志行为）
	if c.Interface != nil {
		c.Interface.Trace(ctx, begin, fc, err)
	}
	sql, _ := fc()
	// 简单判断：以 SELECT 开头的查询计入（忽略大小写与前后空白由上层保证，这里做宽松匹配）
	if len(sql) >= 6 && (sql[:6] == "SELECT" || sql[:6] == "select") {
		c.selects++
	}
}

// TestMemberListAvoidsNPlusOne 验证成员列表批量预取部门/岗位名，而非逐个查。
//
// 历史缺陷：GetUserDepartments / GetDepartmentMembers 在循环内逐个
// dao.GetByID 查部门名、岗位名，N 个成员产生 O(N) 次额外查询（N+1）。
// 现在改为 GetByIDs 一次 IN 查询，查询次数应与成员数量无关。
func TestMemberListAvoidsNPlusOne(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	if _, err := landcorg.InitWithConfig(cfg); err != nil {
		t.Fatalf("InitWithConfig: %v", err)
	}

	// 换装计数 logger，统计后续查询次数。
	db := framedb.GetDB()
	counter := &countingLogger{Interface: db.Logger}
	db.Logger = counter
	// 造数据：3 个部门、3 个岗位、1 个用户挂 3 个部门
	depts := []*model.Department{
		{Name: "研发部", Code: "dept_rd"},
		{Name: "产品部", Code: "dept_pm"},
		{Name: "运营部", Code: "dept_ops"},
	}
	for _, d := range depts {
		if err := db.Create(d).Error; err != nil {
			t.Fatalf("create department: %v", err)
		}
	}
	positions := []*model.Position{
		{Name: "工程师", Code: "pos_eng"},
		{Name: "经理", Code: "pos_mgr"},
		{Name: "专员", Code: "pos_staff"},
	}
	for _, p := range positions {
		if err := db.Create(p).Error; err != nil {
			t.Fatalf("create position: %v", err)
		}
	}

	const userID = 999
	for i := range depts {
		m := &model.OrgMember{
			UserID:       userID,
			DepartmentID: depts[i].ID,
			PositionID:   positions[i].ID,
		}
		if err := db.Create(m).Error; err != nil {
			t.Fatalf("create member: %v", err)
		}
	}

	// 清零计数，聚焦 GetUserDepartments 这一条链路。
	counter.selects = 0
	memberSvc := service.GetMemberService()
	infos, err := memberSvc.GetUserDepartments(userID)
	if err != nil {
		t.Fatalf("GetUserDepartments: %v", err)
	}
	if len(infos) != 3 {
		t.Fatalf("成员部门数 = %d, want 3", len(infos))
	}
	// 部门名、岗位名必须正确组装（批量预取后 map 回填不能串位）。
	wantNames := map[uint]string{}
	for _, d := range depts {
		wantNames[d.ID] = d.Name
	}
	posNames := map[uint]string{}
	for _, p := range positions {
		posNames[p.ID] = p.Name
	}
	for _, info := range infos {
		if info.DepartmentName != wantNames[info.DepartmentID] {
			t.Errorf("DepartmentName = %q, want %q", info.DepartmentName, wantNames[info.DepartmentID])
		}
		if info.PositionName != posNames[info.PositionID] {
			t.Errorf("PositionName = %q, want %q", info.PositionName, posNames[info.PositionID])
		}
	}

	// 关键断言：3 个成员的部门名/岗位名只应来自 2 次批量查询（1 次部门 IN + 1 次岗位 IN），
	// 加上 1 次查成员，共 3 次 SELECT。若仍是 N+1，这里会是 1 + 3 + 3 = 7 次。
	if counter.selects > 3 {
		t.Errorf("GetUserDepartments 产生 %d 次 SELECT，疑似 N+1（期望 ≤3）", counter.selects)
	}
}

// TestDepartmentMembersFillsNames 验证「按部门查成员」也回填部门名与岗位名。
//
// 历史缺陷：GetDepartmentMembers 只预取了岗位名，DepartmentName 一直是空串；
// 而同一文件里的 GetUserDepartments 是有值的。同一个字段在两个接口表现不一致，
// 前端拿到空部门名只能自己再查一次，等于把 N+1 推给了调用方。
func TestDepartmentMembersFillsNames(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	if _, err := landcorg.InitWithConfig(cfg); err != nil {
		t.Fatalf("InitWithConfig: %v", err)
	}

	db := framedb.GetDB()
	counter := &countingLogger{Interface: db.Logger}
	db.Logger = counter

	dept := &model.Department{Name: "技术部", Code: "dept_tech"}
	if err := db.Create(dept).Error; err != nil {
		t.Fatalf("create department: %v", err)
	}
	pos := &model.Position{Name: "后端工程师", Code: "pos_be"}
	if err := db.Create(pos).Error; err != nil {
		t.Fatalf("create position: %v", err)
	}
	for _, uid := range []uint{101, 102, 103} {
		m := &model.OrgMember{UserID: uid, DepartmentID: dept.ID, PositionID: pos.ID}
		if err := db.Create(m).Error; err != nil {
			t.Fatalf("create member: %v", err)
		}
	}

	counter.selects = 0
	infos, total, err := service.GetMemberService().GetDepartmentMembers(dept.ID, 1, 20)
	if err != nil {
		t.Fatalf("GetDepartmentMembers: %v", err)
	}
	if total != 3 || len(infos) != 3 {
		t.Fatalf("成员数 = %d (total=%d), want 3", len(infos), total)
	}

	for _, info := range infos {
		if info.DepartmentName != "技术部" {
			t.Errorf("user %d 的 DepartmentName = %q, want 技术部", info.UserID, info.DepartmentName)
		}
		if info.PositionName != "后端工程师" {
			t.Errorf("user %d 的 PositionName = %q, want 后端工程师", info.UserID, info.PositionName)
		}
	}

	// 3 个成员应只产生常数次查询（count + 成员 + 部门 IN + 岗位 IN），
	// 若退回 N+1 会随成员数线性增长。
	if counter.selects > 5 {
		t.Errorf("GetDepartmentMembers 产生 %d 次 SELECT，疑似 N+1（期望 ≤5）", counter.selects)
	}
}
