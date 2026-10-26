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
//
// service_auth.key 是本系统鉴权的服务间调用凭证：测试用它证明「内部服务能进、
// 外部匿名不能进」。jwt.secret 用固定测试值。
const deployConfigTemplate = `
server:
  addr: "127.0.0.1"
  port: 0
  use_default_routes: true
database:
  driver: "sqlite"
  dsn: "%s"
  auto_migrate: true
jwt:
  secret: "test-only-secret-at-least-32-characters-long"
  issuer: "landc-org-test"
log:
  level: "error"
  output: "console"
other:
  server:
    gin_mode: "test"
  service_auth:
    key: "%s"
`

// testServiceKey 是测试用的服务间凭证。
const testServiceKey = "test-service-key-for-e2e"

// postAsService 以内部服务身份发起请求（带 X-Service-Key）。
func postAsService(t *testing.T, url, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", testServiceKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	return resp
}

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

	if err := os.WriteFile(cfgPath, []byte(fmt.Sprintf(deployConfigTemplate, filepath.ToSlash(dbPath), testServiceKey)), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	// InitDefault 内部走 InitGlobalConfigWithPath（once 语义），若全局配置已被前序
	// 测试占用则不会重载，导致连到已删除的旧 DB。先显式重置到本测试的配置。
	resetFramework(t, cfgPath)
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

	resp := postAsService(t, srv.URL+"/api/department/create", `{"name":"研发中心","code":"rd"}`)
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
	if err := os.WriteFile(cfgPath, []byte(fmt.Sprintf(deployConfigTemplate, filepath.ToSlash(dbPath), testServiceKey)), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	resetFramework(t, cfgPath)
	t.Cleanup(func() { _ = framedb.Close() })

	server, err := landcorg.InitDefault(cfgPath)
	if err != nil {
		t.Fatalf("InitDefault: %v", err)
	}
	srv := httptest.NewServer(server.Engine())
	defer srv.Close()

	resp := postAsService(t, srv.URL+"/api/department/create", `{}`)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var count int64
		_ = framedb.GetDB().Model(&landcorgModel.Department{}).Count(&count).Error
		t.Fatalf("空请求体应被拒绝，实际 HTTP %d，响应 %s，部门记录数 %d", resp.StatusCode, raw, count)
	}
}

// TestUnauthenticatedRejected 守住「本系统必须有鉴权」这个不变量。
//
// 历史缺陷：InitWithConfig 曾只挂 trace 与日志中间件，漏了鉴权 —— 任何人只要
// 知道地址就能创建部门、添加成员、分配组织角色。同为内网服务的
// LandcPermission 一直挂着框架 middleware.Auth()，本系统属于遗漏。
//
// 本用例无凭证请求，必须 401；同时用健康检查证明探针不受影响。
func TestUnauthenticatedRejected(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "data", "org.db")
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(fmt.Sprintf(deployConfigTemplate, filepath.ToSlash(dbPath), testServiceKey)), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	resetFramework(t, cfgPath)
	t.Cleanup(func() { _ = framedb.Close() })

	server, err := landcorg.InitDefault(cfgPath)
	if err != nil {
		t.Fatalf("InitDefault: %v", err)
	}
	srv := httptest.NewServer(server.Engine())
	defer srv.Close()

	// 1. 无任何凭证的业务请求必须被拒绝
	resp, err := http.Post(srv.URL+"/api/department/create", "application/json",
		strings.NewReader(`{"name":"越权部门","code":"hacked"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("无凭证请求应返回 401，实际 HTTP %d，响应 %s", resp.StatusCode, raw)
	}

	// 2. 错误的服务密钥同样必须被拒绝（不能只看头是否存在）
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/department/create",
		strings.NewReader(`{"name":"越权部门","code":"hacked2"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", "wrong-key")
	badResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post with wrong key: %v", err)
	}
	defer badResp.Body.Close()
	if badResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("错误服务密钥应返回 401，实际 HTTP %d", badResp.StatusCode)
	}

	// 3. 健康检查不受鉴权影响（否则 k8s 探针永远失败）
	healthResp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	defer healthResp.Body.Close()
	if healthResp.StatusCode != http.StatusOK {
		t.Fatalf("健康检查应返回 200，实际 HTTP %d", healthResp.StatusCode)
	}

	// 4. 越权请求不能落库
	var count int64
	_ = framedb.GetDB().Model(&landcorgModel.Department{}).Count(&count).Error
	if count != 0 {
		t.Fatalf("被鉴权拦截的请求不应写入数据，实际部门记录数 %d", count)
	}
}
