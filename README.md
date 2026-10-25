# LandcOrg

组织架构管理系统，提供部门管理、岗位管理、成员管理、邀请加入、组织驱动授权等功能，支持树形组织结构和灵活的角色继承策略。

## 功能特性

### 部门管理
- 部门 CRUD（树形结构，支持无限层级）
- 部门树查询（递归构建完整组织树）
- 自动计算层级（level）和路径（path）
- 部门排序

### 岗位管理
- 岗位/职级 CRUD
- 岗位等级划分

### 成员管理
- 添加/移除成员
- 支持一人多部门
- 部门负责人标记
- 成员岗位设置

### 邀请机制
- 生成邀请码（支持指定目标用户或开放邀请）
- 接受邀请自动加入部门
- 邀请过期时间控制
- 邀请使用次数限制（单次/多次/无限）
- 撤销邀请

### 组织驱动授权
- 部门-角色绑定（加入部门自动继承角色）
- 岗位-角色绑定（担任岗位自动继承角色）
- 查询用户从组织结构继承的所有角色
- 与 LandcPermission 联动实现完整 RBAC

## 技术栈

- **框架**：`github.com/LandcLi/landc-go/frame`（配置 / DB / 缓存 / JWT / DI / Web / 健康检查）
- **Web**：Gin
- **数据库**：SQLite（默认，纯 Go 驱动，零外部依赖）或 MySQL（生产）
- **缓存**：Redis（可选；缺失时降级为进程内 LRU）
- **日志**：`landc-go/log`

## 项目结构

```
LandcOrg/
├── api/                        接口定义 + 路由表（department / position / member / invitation / org_role）
├── dao/                        数据访问接口
├── internal/                   内部实现（可通过 DI Override 替换）
│   ├── controller/
│   ├── dao_impl/
│   └── service_impl/
├── model/                      数据模型（部门 / 岗位 / 成员 / 邀请 / 角色绑定）
├── sdk/                        HTTP SDK（远程调用）
├── service/                    服务接口
├── sqls/                       MySQL 初始化脚本
├── tests/sqlite/               SQLite 端到端冒烟
├── utility/                    配置 / 日志（DB 与缓存委托给框架）
├── cmd/api/main.go             独立运行入口
└── LandcOrg.go                 库入口（InitDefault / InitWithConfig / InitWithConfigAndGin）
```

## 快速开始

### 1. 准备配置并启动

```bash
cp config.example.yml config.yaml
go run ./cmd/api
```

默认使用 SQLite，首次启动自动在 `data/landc_org.db` 建表
（`database.auto_migrate: true`），无需 MySQL / Redis。

```bash
curl http://127.0.0.1:8083/health
```

创建部门：

```bash
curl -X POST http://127.0.0.1:8083/api/department/create \
  -H 'Content-Type: application/json' \
  -d '{"name":"研发中心","code":"rd","description":"默认部门"}'
```

### 2. 切换到 MySQL（生产）

```yaml
database:
  driver: "mysql"
  dsn: "root:password@tcp(127.0.0.1:3306)/landc_org?charset=utf8mb4&parseTime=True&loc=Local"
  auto_migrate: false
```

```bash
mysql -u root -p < sqls/init.sql
```

### 3. 嵌入到其他系统（库模式）

```go
// 方式一：使用独立 Server（含中间件与路由）
server, err := landcorg.InitDefault("config.yaml")

// 方式二：挂到宿主已有的 gin 引擎（不自带中间件链，由宿主决定）
landcorg.InitWithConfigAndGin(orgCfg, ginEngine)
```

库模式只注册控制器，不会创建 engine、不会调用 `Run`，宿主可继续追加自有路由。

### 4. Docker 部署

```bash
docker build -t landc-org:0.1.1 .

docker run -d --name landc-org \
  -p 8083:8083 \
  -v "$(pwd)/data:/app/data" \
  landc-org:0.1.1
```

镜像要点：

- 多阶段构建，产物是 `CGO_ENABLED=0` 静态二进制（纯 Go SQLite 驱动，运行阶段无需 gcc）
- 以非 root 用户 `app` 运行，`/app/data` 与 `/app/logs` 已预建并授权
- 内置 `HEALTHCHECK` 打 `/health`，可直接被 compose / k8s 探活
- 时区固定 `Asia/Shanghai`（`tzdata` 已安装）

挂载 `data/` 是关键：不挂载时容器重建会丢掉 SQLite 库。


## 使用 SDK（远程调用模式）

```go
import (
    "github.com/LandcLi/LandcOrg/api/department"
    "github.com/LandcLi/LandcOrg/api/invitation"
    "github.com/LandcLi/LandcOrg/sdk"
)

// 注册远程实现
department.DepartmentGateway.Override(
    sdk.NewRemoteDepartmentController("http://org-service:8083"),
)
invitation.InvitationGateway.Override(
    sdk.NewRemoteInvitationController("http://org-service:8083"),
)

// 使用方式与本地完全一致
ctrl := department.GetDepartmentController()
resp, err := ctrl.GetDepartmentTree(ctx, &v1.GetDepartmentTreeRequest{RootID: 0})
```

## 业务流程

### 组织创建与邀请加入

```
1. 管理员创建部门
   POST /api/department/create
   → 创建根部门或子部门

2. 管理员设置部门默认角色
   POST /api/org-role/assign-department-roles
   → 加入该部门的人自动获得这些角色

3. 管理员生成邀请码
   POST /api/invitation/create
   → 指定目标部门、岗位、过期时间

4. 被邀请人接受邀请
   POST /api/invitation/accept
   → 自动加入部门 + 设置岗位
   → 自动继承部门角色 + 岗位角色
```

### 权限校验流程（与 LandcPermission 联动）

```
请求到达 → Token 校验 → 获取 userID
→ LandcPermission.CheckPermission(userID, resource, action)
  → 获取直接分配的角色
  → 获取组织继承的角色（调用 LandcOrg.GetUserInheritedRoles）
  → 合并所有角色的权限
  → 匹配 resource + action
  → 返回 allowed: true/false
```

## API 接口

### 部门接口

| 接口 | 路径 | 描述 |
|------|------|------|
| CreateDepartment | POST /api/department/create | 创建部门 |
| UpdateDepartment | POST /api/department/update | 更新部门 |
| DeleteDepartment | POST /api/department/delete | 删除部门（无子部门时） |
| GetDepartment | POST /api/department/get | 获取部门详情 |
| GetDepartmentTree | POST /api/department/tree | 获取部门树 |
| ListDepartments | POST /api/department/list | 部门列表（分页） |

### 岗位接口

| 接口 | 路径 | 描述 |
|------|------|------|
| CreatePosition | POST /api/position/create | 创建岗位 |
| UpdatePosition | POST /api/position/update | 更新岗位 |
| DeletePosition | POST /api/position/delete | 删除岗位 |
| GetPosition | POST /api/position/get | 获取岗位详情 |
| ListPositions | POST /api/position/list | 岗位列表（分页） |

### 成员接口

| 接口 | 路径 | 描述 |
|------|------|------|
| AddMember | POST /api/member/add | 添加成员到部门 |
| RemoveMember | POST /api/member/remove | 从部门移除成员 |
| UpdateMember | POST /api/member/update | 更新成员岗位/负责人状态 |
| GetUserDepartments | POST /api/member/get-user-departments | 获取用户所属部门 |
| GetDepartmentMembers | POST /api/member/get-department-members | 获取部门成员列表 |

### 邀请接口

| 接口 | 路径 | 描述 |
|------|------|------|
| CreateInvitation | POST /api/invitation/create | 创建邀请 |
| AcceptInvitation | POST /api/invitation/accept | 接受邀请 |
| RevokeInvitation | POST /api/invitation/revoke | 撤销邀请 |
| ListInvitations | POST /api/invitation/list | 我的邀请列表 |

### 组织-角色绑定接口

| 接口 | 路径 | 描述 |
|------|------|------|
| AssignDepartmentRoles | POST /api/org-role/assign-department-roles | 绑定部门默认角色 |
| RevokeDepartmentRoles | POST /api/org-role/revoke-department-roles | 解绑部门角色 |
| GetDepartmentRoles | POST /api/org-role/get-department-roles | 获取部门绑定的角色 |
| AssignPositionRoles | POST /api/org-role/assign-position-roles | 绑定岗位默认角色 |
| RevokePositionRoles | POST /api/org-role/revoke-position-roles | 解绑岗位角色 |
| GetPositionRoles | POST /api/org-role/get-position-roles | 获取岗位绑定的角色 |
| GetUserInheritedRoles | POST /api/org-role/get-user-inherited-roles | 获取用户继承的所有角色 |

## 数据模型

### 部门树结构

```
总公司 (parent_id=0, level=1, path="/")
├── 技术部 (parent_id=1, level=2, path="/1/")
│   ├── 后端组 (parent_id=2, level=3, path="/1/2/")
│   └── 前端组 (parent_id=2, level=3, path="/1/2/")
├── 产品部 (parent_id=1, level=2, path="/1/")
└── 运营部 (parent_id=1, level=2, path="/1/")
```

### 授权继承模型

```
用户张三
├── 所属部门：技术部 → 绑定角色：developer
├── 所属部门：后端组 → 绑定角色：backend_dev
├── 担任岗位：高级工程师 → 绑定角色：senior
└── 直接分配角色：project_lead

张三的有效角色 = {developer, backend_dev, senior, project_lead}
```

## 自定义扩展

所有接口层都可以通过 DI Override 替换：

```go
// 替换部门服务逻辑
service.RegisterDepartmentService(myCustomDeptService)

// 替换成员 DAO（如使用 MongoDB）
dao.RegisterMemberDAO(myMongoMemberDAO)

// 替换邀请服务（如增加审批流程）
service.RegisterInvitationService(myApprovalInvitationService)
```

## 与其他项目的关系

```
LandcSSO          ← 认证授权协议（Token/OAuth）
LandcUserManager  ← 用户身份数据（CRUD）
LandcPermission   ← 权限管理（RBAC）
LandcOrg          ← 组织架构（部门/岗位/成员/邀请/角色继承）
landc-go/mvc      ← 底层框架（DI/Gateway/Web）
```

## 许可证

MIT License
