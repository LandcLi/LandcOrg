package sqlite_test

import (
	"context"
	"net/http/httptest"
	"testing"

	framedb "github.com/LandcLi/landc-go/frame/pkg/db"
	"github.com/gin-gonic/gin"

	landcorg "github.com/LandcLi/LandcOrg"
	deptV1 "github.com/LandcLi/LandcOrg/api/department/v1"
	orgRoleV1 "github.com/LandcLi/LandcOrg/api/org_role/v1"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/sdk"
)

// TestSDKRouteContract 验证 SDK 的远程调用路径与服务端注册路径一致。
//
// 框架的 ProxyClient 从请求结构体的 meta.Meta 推导 path/method，缺失时回退到
// "/"+小写方法名。历史版本 Org 的请求结构体没有 meta 标签，因此 SDK 全部方法
// 都会打到不存在的路径（如 /createdepartment）返回 404 —— 而服务端测试完全
// 覆盖不到这一点：服务端用自带的 XxxRoutes 路由表注册，自身链路是通的。
//
// 这个测试是「服务端与 SDK 共用一份 meta 声明」这条约束的守门人：
// 任何人给请求结构体改错 path，或删掉 meta 标签，这里都会红。
func TestSDKRouteContract(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	engine := gin.New()
	if err := landcorg.InitWithConfigAndGin(cfg, engine); err != nil {
		t.Fatalf("InitWithConfigAndGin: %v", err)
	}
	srv := httptest.NewServer(engine)
	defer srv.Close()

	ctx := context.Background()

	// 1. 通过 SDK 建部门：path 或 method 不对都会拿不到 200
	deptClient := sdk.NewRemoteDepartmentController(srv.URL)
	created, err := deptClient.CreateDepartment(ctx, &deptV1.CreateDepartmentRequest{
		Name: "SDK 部门",
		Code: "sdk-dept",
	})
	if err != nil {
		t.Fatalf("sdk CreateDepartment: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("sdk CreateDepartment returned zero id: %+v", created)
	}

	// 2. 通过 SDK 查用户继承角色：路径正确，且未加入组织时结果为空
	roleClient := sdk.NewRemoteOrgRoleController(srv.URL)
	inherited, err := roleClient.GetUserInheritedRoles(ctx, &orgRoleV1.GetUserInheritedRolesRequest{UserID: 1})
	if err != nil {
		t.Fatalf("sdk GetUserInheritedRoles: %v", err)
	}
	if len(inherited.RoleIDs) != 0 {
		t.Errorf("user without org membership should inherit no roles, got %v", inherited.RoleIDs)
	}

	// 3. 确认请求体真的被服务端解析并落库（而不是空 body 蒙混过关）
	var count int64
	if err := framedb.GetDB().Model(&model.Department{}).Where("code = ?", "sdk-dept").Count(&count).Error; err != nil {
		t.Fatalf("count department: %v", err)
	}
	if count != 1 {
		t.Errorf("department created via sdk not persisted, count=%d", count)
	}
}
