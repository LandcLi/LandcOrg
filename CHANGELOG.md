# Changelog

本文件记录 LandcOrg 的显著变更。
格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [0.1.2] - 2026-10-26

### 修复

修复三个会**实质损坏组织数据**的缺陷：

- **移动部门不重算子树层级/路径**。移动只改 `parent_id`，`level`/`path` 停在旧值，
  按层级做数据权限、按路径做前缀匹配的功能全部算错。现在移动时连同整棵子树一起重算。
- **允许部门挂到自己或子孙下（成环）**。成环后 `GetDepartmentTree` 父子互相引用，
  json 序列化因循环引用失败，整棵树读不出来。现在沿目标父级祖先链校验拒绝成环；
  对库里已存在的环则降级为「环上节点当根」，接口永远返回合法树。
- **删除部门不拦成员、不清角色绑定**。成员变成孤儿（「我的部门」查不到）、
  已删部门的角色继续算给用户。现在有成员拒绝删除，删除时一并清理 `department_role`。

### 测试

- `tests/sqlite/department_consistency_test.go`：5 个回归用例，覆盖移动重算、
  拒绝成环、拦成员、清角色绑定，以及直接写库伪造环数据的防御性用例。

## [0.1.1] - 2026-10-25

修复三个「服务能起来、接口全哑」的严重缺陷，并升级到 `frame v0.8.6` / `api v0.5.1`。

### 修复

- **SDK 调用全线 404（严重）**。27 个请求结构体缺少 `meta.Meta` 路由声明，
  而 `di.ProxyClient` 与 `web.RegisterHandler` 都只从该标签反射 path/method；
  缺失时客户端回退到 `"/" + 小写方法名`，与真实路由（`/api/department/create`）不符。
  表现为 SDK 调用必 404，服务端却看不到任何日志。现已为全部请求结构体补齐标签。
- **独立部署时业务路由完全不注册（严重）**。`utility/config` 曾用私有
  `gopkg.in/yaml.v2` 结构二次解析同一份配置，字段名写成了 `use_default_route`（单数），
  而框架字段是 `use_default_routes`（复数）。重复解析后该值恒为 `false`，
  业务路由静默不注册 —— `/health` 正常、任意 `/api/*` 都落到 NoRoute。
  现改为**配置单一解析源**：`config.LoadConfig` 委托框架加载，`config.FromFramework`
  只做纯映射，不再自持 yaml 结构。
- **`/api/*` 错误响应统一返回 HTTP 200**：随 `api v0.5.1` 修复（框架日志中间件
  吞掉状态码），参数校验失败现在正确返回 400。

### 变更

- **路由声明收敛为单一来源**：删除 5 张 `var XxxRoutes = di.Routes{...}` 路由表，
  路由信息只存在于请求结构体的 `meta.Meta` 标签，服务端注册与 SDK 共用同一份声明，
  杜绝「加了接口忘了同步路由表」的双写漂移。
- `utility/config` 里的自定义 `other` 下钻逻辑替换为框架的
  `GetOtherString` / `GetOtherBool`，`gopkg.in/yaml.v2` 从直接依赖中移除。

### 测试

- `tests/sqlite/sdk_contract_test.go`：用真实 `sdk.NewRemoteDepartmentController`
  对 httptest 服务端往返，断言 path / method / 落库结果。此前 SDK 与路由表的
  契约无任何测试覆盖，是 404 缺陷长期未被发现的原因。
- `tests/sqlite/deploy_path_test.go`：按**部署配置模板**（框架字段名）走一遍
  `InitDefault`，断言 27 条业务路由注册成功、空请求体返回 400 且不落库。

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
