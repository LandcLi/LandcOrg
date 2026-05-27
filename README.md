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

- **框架**: landc-go/mvc（DI / Gateway / Web）
- **Web框架**: Gin
- **数据库**: MySQL (GORM)
- **缓存**: Redis
- **日志**: landc-go/log
- **API规范**: landc-go/api

## 项目结构

```
LandcOrg/
├── api/                        # API接口定义
│   ├── department/             # 部门管理接口
│   ├── position/               # 岗位管理接口
│   ├── member/                 # 成员管理接口
│   ├── invitation/             # 邀请管理接口
│   └── org_role/               # 组织-角色绑定接口
├── dao/                        # 数据访问层接口
├── internal/                   # 内部实现（可通过 DI Override 替换）
│   ├── controller/             # 控制器层
│   ├── dao_impl/              # DAO实现层
│   └── service_impl/          # 服务实现层
├── model/                      # 数据模型
├── sdk/                        # HTTP SDK（远程调用）
├── service/                    # 服务层接口
├── sqls/                       # SQL脚本
├── utility/                    # 工具包
│   ├── config/                # 配置管理
│   ├── db/                    # 数据库连接
│   └── logger/                # 日志管理
├── cmd/api/main.go            # 独立运行入口
└── LandcOrg.go                 # 主入口文件（支持嵌入集成）
```

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p < sqls/init.sql
```

### 2. 配置文件

```bash
cp config.example.yml config.yml
# 修改数据库、Redis等配置
```

### 3. 独立运行

```go
package main

import (
    "fmt"
    "log"
    landcorg "github.com/LandcLi/LandcOrg"
)

func main() {
    engine, err := landcorg.InitDefault("config.yml")
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    fmt.Println("LandcOrg starting on :8083")
    engine.Run(":8083")
}
```

### 4. 嵌入到其他系统

```go
// 本地集成模式
orgCfg := &config.Config{...}
landcorg.InitWithConfigAndGin(orgCfg, ginEngine)
```

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
