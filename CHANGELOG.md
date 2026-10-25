# Changelog

本文件记录 LandcOrg 的显著变更。
格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [0.1.0] - 2026-10-25

首个可发布版本。

### 新增

- **SQLite 开箱即用**：默认 `driver: sqlite`，纯 Go 驱动（`modernc.org/sqlite`），
  不依赖 CGO / gcc / 外部数据库服务。`database.auto_migrate: true` 时首次启动自动建表。
- **`Models()` / `EnsureSchema()`**：模型清单与建表入口唯一来源，供宿主与测试复用。
- **两种库模式**：`InitWithConfigAndGin`（挂到宿主 gin 引擎）与
  `RegisterHandlers`（只注册控制器，中间件链由宿主决定）。
- **`config.FacadeLogConfig()`**：把本系统配置的 log 段转换为框架 `facade.LogConfig`，
  不再自持一份 `Log` 字段（与框架的 string 级别枚举类型冲突，会让 YAML 解析必然失败）。
- **CI**（`.github/workflows/ci.yml`）：vet / gofmt / `go test -race` / 静态二进制构建 /
  镜像构建校验。
- **Dockerfile**：多阶段构建、非 root 运行、内置 `/health` HEALTHCHECK、固定时区。
- 冒烟测试 `tests/sqlite`：建表、库模式与宿主路由共存、HTTP 创建部门。

### 变更

- **基础设施统一委托框架**：`utility/db` 从自建 gorm + go-redis v8 改为薄委托
  （`GetDB()` → `framedb.GetDB()`、`GetRedis()` → `framecache.GetRedis()`），
  消除与框架重复的连接管理与配置解析。
- **路由注册改用 `web.RegisterHandlers` + `di.Routes` 注册选项**：
  框架已移除 `Gateway.RegisterRoutes`，改为复用既有的
  `di.Routes` 路由表生成 `web.WithMethodPath` / `WithMethodHTTPMethod` 选项，
  避免给数十个请求结构体逐个补 `meta` 标签。
- `di.NewGateway[T]` / `di.NewProxyClient` 去掉多余的路由表参数（框架签名变更）。
- `cmd/api/main.go`：支持 `LANDC_CONFIG` 环境变量，并用 `signal.NotifyContext` 优雅停机。
- `config.example.yml`：`server.port` 由字符串 `":8083"` 改为
  `addr` + `port: 8083`（int 字段解析不了带冒号的字符串）；
  补 `database.driver`，否则框架判定为「未配置数据库」直接跳过连接，运行期拿到 nil DB。
- 依赖升级到 `landc-go/frame v0.8.5`、`landc-go/api & log v0.5.0`。

### 修复

- **DAO 缓存永久失效**：各 DAO 在包 `init()` 中 `db.GetRedis()` 并把结果存进结构体字段，
  而 Redis 要等 `main` 之后才初始化，字段恒为 `nil`，缓存从不命中且不报错。
  改为调用时惰性获取。
- `.gitattributes` 统一 LF，消除 `core.autocrlf=true` 环境下 `gofmt -l` 的全量误报。
