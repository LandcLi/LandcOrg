package sqlite_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	framedb "github.com/LandcLi/landc-go/frame/pkg/db"

	landcorg "github.com/LandcLi/LandcOrg"
	landcorgModel "github.com/LandcLi/LandcOrg/model"
)

// deployConfigTemplate 按框架文档的字段名书写（server.use_default_routes），
// 模拟用户照着 config.example.yml 配置的真实情况。
const deployConfigTemplate = `
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

// TestInitDefaultRegistersBusinessRoutes 锁定「独立部署路径」的行为。
//
// 历史缺陷：InitDefault 曾用本系统自己的 yaml 结构二次解析配置文件，字段名
// 与框架不一致（本系统 use_default_route，框架 use_default_routes）。用户按
// 框架文档写配置时，本系统读到 false → 业务控制器整体不注册 → 所有 /api/*
// 请求落到 NoRoute，返回 200 + 空 data。接口看起来「通了」，实际什么都没做。
//
// 此前所有测试都走 InitWithConfigAndGin（库模式），完全覆盖不到真实部署路径，
// 所以这个缺陷能一直潜伏。本用例是该路径的守门人。
func TestInitDefaultRegistersBusinessRoutes(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data", "org.db")
	cfgPath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(cfgPath, []byte(fmt.Sprintf(deployConfigTemplate, filepath.ToSlash(dbPath))), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Cleanup(func() { _ = framedb.Close() })

	server, err := landcorg.InitDefault(cfgPath)
	if err != nil {
		t.Fatalf("InitDefault: %v", err)
	}

	// 1. 路由表必须包含业务路由（而不是只有健康检查）
	want := "POST /api/department/create"
	found := false
	for _, r := range server.Routes() {
		if r.Method+" "+r.Path == want {
			found = true
			break
		}
	}
	if !found {
		paths := make([]string, 0, len(server.Routes()))
		for _, r := range server.Routes() {
			paths = append(paths, r.Method+" "+r.Path)
		}
		t.Fatalf("业务路由未注册，实际路由: %v", paths)
	}

	// 2. 真实请求必须命中 handler 并带回 data
	srv := httptest.NewServer(server.Engine())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/department/create", "application/json",
		strings.NewReader(`{"name":"研发中心","code":"rd"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var envelope struct {
		Code int `json:"code"`
		Data *struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("响应不是合法 JSON: %s", raw)
	}
	if envelope.Data == nil {
		t.Fatalf("响应缺少 data 字段（说明请求落到了 NoRoute 而非业务 handler）: %s", raw)
	}
	if envelope.Data.ID == 0 {
		t.Fatalf("data.id 为 0，创建未生效: %s", raw)
	}
	if envelope.Data.Code != "rd" {
		t.Errorf("data.code = %q, want rd", envelope.Data.Code)
	}
}

// TestInvalidRequestRejected 验证缺必填字段的请求被拒绝。
//
// 若 binding:"required" 不生效，空请求体会被一路放行到 service 层，
// 创建出 name/code 均为空串的脏数据，且客户端拿到 200 无从察觉。
func TestInvalidRequestRejected(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data", "org.db")
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(fmt.Sprintf(deployConfigTemplate, filepath.ToSlash(dbPath))), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Cleanup(func() { _ = framedb.Close() })

	server, err := landcorg.InitDefault(cfgPath)
	if err != nil {
		t.Fatalf("InitDefault: %v", err)
	}
	srv := httptest.NewServer(server.Engine())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/department/create", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var count int64
		_ = framedb.GetDB().Model(&landcorgModel.Department{}).Count(&count).Error
		t.Fatalf("空请求体应被拒绝，实际 HTTP %d，响应 %s，部门记录数 %d", resp.StatusCode, raw, count)
	}
}
