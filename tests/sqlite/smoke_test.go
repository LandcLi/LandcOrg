package sqlite_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	framedb "github.com/LandcLi/landc-go/frame/pkg/db"
	"github.com/gin-gonic/gin"

	landcorg "github.com/LandcLi/LandcOrg"
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/config"
)

const configTemplate = `
server:
  addr: "127.0.0.1"
  port: 0
  use_default_routes: true
database:
  driver: "sqlite"
  dsn: "%s"
  auto_migrate: true
log:
  level: "error"
  output: "console"
other:
  server:
    gin_mode: "test"
`

// setupSQLiteEnv 用临时目录中的 SQLite 完成一次完整初始化，返回 Org 配置。
func setupSQLiteEnv(t *testing.T) *config.Config {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data", "org.db")
	configPath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(fmt.Sprintf(configTemplate, filepath.ToSlash(dbPath))), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	resetFramework(t, configPath)
	t.Cleanup(func() { _ = framedb.Close() })

	return config.FromFramework(frameconfig.GetConfig())
}

// resetFramework 把框架全局配置与全局 DB 重置到指定配置文件。
//
// 关键：frameconfig.InitGlobalConfigWithPath 与 framedb.InitGlobalDBWithDefault
// 都是 once 语义（已初始化则跳过），而测试要各自用独立 SQLite 文件。若不显式
// 重置，第二个测试会继续读到第一个测试的（已随 t.TempDir 删除的）DB 文件 → 500。
// 这里先关闭全局 DB（置 nil），再用 InitGlobalConfigWithConfig（无条件覆盖）切到新配置。
func resetFramework(t *testing.T, configPath string) {
	t.Helper()
	_ = framedb.Close()

	fc, err := frameconfig.LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("load framework config: %v", err)
	}
	if err := frameconfig.InitGlobalConfigWithConfig(fc); err != nil {
		t.Fatalf("init framework config: %v", err)
	}
}

// TestSQLiteSchemaCreated 验证 auto_migrate 建出了组织域六张表。
func TestSQLiteSchemaCreated(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	if _, err := landcorg.InitWithConfig(cfg); err != nil {
		t.Fatalf("InitWithConfig: %v", err)
	}

	migrator := framedb.GetDB().Migrator()
	for _, m := range landcorg.Models() {
		if !migrator.HasTable(m) {
			t.Errorf("table missing after EnsureSchema: %T", m)
		}
	}
}

// TestLibraryModeHostRoutes 验证库模式：宿主已有路由不受影响。
func TestLibraryModeHostRoutes(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	engine := gin.New()
	if err := landcorg.InitWithConfigAndGin(cfg, engine); err != nil {
		t.Fatalf("InitWithConfigAndGin: %v", err)
	}

	engine.GET("/host/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	srv := httptest.NewServer(engine)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/host/ping")
	if err != nil {
		t.Fatalf("get host route: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("host route status = %d, want 200", resp.StatusCode)
	}
}

// TestCreateDepartmentOverHTTP 走完整 HTTP 链路创建部门，验证路由映射与 SQLite 落库。
func TestCreateDepartmentOverHTTP(t *testing.T) {
	cfg := setupSQLiteEnv(t)

	engine := gin.New()
	if err := landcorg.InitWithConfigAndGin(cfg, engine); err != nil {
		t.Fatalf("InitWithConfigAndGin: %v", err)
	}
	srv := httptest.NewServer(engine)
	defer srv.Close()

	body, _ := json.Marshal(map[string]interface{}{
		"name":        "研发中心",
		"code":        "rd",
		"description": "sqlite smoke",
	})
	resp, err := http.Post(srv.URL+"/api/department/create", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post create department: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create department status = %d, raw=%s", resp.StatusCode, string(raw))
	}

	var count int64
	if err := framedb.GetDB().Model(&model.Department{}).Where("code = ?", "rd").Count(&count).Error; err != nil {
		t.Fatalf("count departments: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 department persisted, got %d (raw=%s)", count, string(raw))
	}
}
