package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/LandcLi/landc-go/frame/pkg/web"

	"github.com/LandcLi/LandcOrg/api/department"
	"github.com/LandcLi/LandcOrg/api/invitation"
	"github.com/LandcLi/LandcOrg/api/member"
	"github.com/LandcLi/LandcOrg/api/org_role"
	"github.com/LandcLi/LandcOrg/api/position"
)

// 路由声明的唯一来源是各请求结构体上的 meta.Meta 标签。
//
// 服务端注册（web.RegisterHandler）与远程 SDK（di.ProxyClient）读取的是
// 同一份声明：前者反射出 path/method 生成路由，后者反射出 path/method 拼请求。
// 因此不存在「服务端与 SDK 各写一份路由、改一处漏一处」的漂移可能。
//
// 历史上本项目用 api/<模块>/XxxRoutes 变量集中声明路由，并通过注册选项
// 覆盖到框架上。那种做法让 SDK 的 ProxyClient 拿不到路由（它只认
// meta.Meta，缺失时回退到 "/+小写方法名"），导致远程调用全部 404。

// RegisterDefaultRoutes 把本系统全部控制器注册到 server（独立部署使用）。
// 健康检查由框架按 server.use_default_routes 自动注册。
func RegisterDefaultRoutes(server *web.Server) error {
	registrations := []struct {
		name string
		fn   func() error
	}{
		{"department", func() error { return server.RegisterHandler(department.GetDepartmentController()) }},
		{"position", func() error { return server.RegisterHandler(position.GetPositionController()) }},
		{"member", func() error { return server.RegisterHandler(member.GetMemberController()) }},
		{"invitation", func() error { return server.RegisterHandler(invitation.GetInvitationController()) }},
		{"org_role", func() error { return server.RegisterHandler(org_role.GetOrgRoleController()) }},
	}

	for _, r := range registrations {
		if err := r.fn(); err != nil {
			return fmt.Errorf("注册 %s 控制器失败: %w", r.name, err)
		}
	}
	return nil
}

// RegisterHandlers 把本系统全部控制器注册到任意 gin 路由（嵌入到宿主时使用）。
// 不挂载全局中间件：中间件链由宿主决定。
func RegisterHandlers(router gin.IRouter) error {
	registrations := []struct {
		name string
		fn   func() error
	}{
		{"department", func() error {
			_, err := web.RegisterHandlers(router, department.GetDepartmentController())
			return err
		}},
		{"position", func() error {
			_, err := web.RegisterHandlers(router, position.GetPositionController())
			return err
		}},
		{"member", func() error {
			_, err := web.RegisterHandlers(router, member.GetMemberController())
			return err
		}},
		{"invitation", func() error {
			_, err := web.RegisterHandlers(router, invitation.GetInvitationController())
			return err
		}},
		{"org_role", func() error {
			_, err := web.RegisterHandlers(router, org_role.GetOrgRoleController())
			return err
		}},
	}

	for _, r := range registrations {
		if err := r.fn(); err != nil {
			return fmt.Errorf("注册 %s 控制器失败: %w", r.name, err)
		}
	}
	return nil
}
