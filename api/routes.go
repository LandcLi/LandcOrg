package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/LandcLi/landc-go/frame/pkg/di"
	"github.com/LandcLi/landc-go/frame/pkg/web"

	"github.com/LandcLi/LandcOrg/api/department"
	"github.com/LandcLi/LandcOrg/api/invitation"
	"github.com/LandcLi/LandcOrg/api/member"
	"github.com/LandcLi/LandcOrg/api/org_role"
	"github.com/LandcLi/LandcOrg/api/position"
)

// RegisterDefaultRoutes 把本系统全部控制器注册到 server（独立部署使用）。
// 健康检查由框架按 server.use_default_routes 自动注册。
func RegisterDefaultRoutes(server *web.Server) error {
	registrations := []struct {
		name string
		fn   func() error
	}{
		{"department", func() error {
			return server.RegisterHandler(department.GetDepartmentController(), optionsFrom(department.DepartmentRoutes)...)
		}},
		{"position", func() error {
			return server.RegisterHandler(position.GetPositionController(), optionsFrom(position.PositionRoutes)...)
		}},
		{"member", func() error {
			return server.RegisterHandler(member.GetMemberController(), optionsFrom(member.MemberRoutes)...)
		}},
		{"invitation", func() error {
			return server.RegisterHandler(invitation.GetInvitationController(), optionsFrom(invitation.InvitationRoutes)...)
		}},
		{"org_role", func() error {
			return server.RegisterHandler(org_role.GetOrgRoleController(), optionsFrom(org_role.OrgRoleRoutes)...)
		}},
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
			_, err := web.RegisterHandlers(router, department.GetDepartmentController(), optionsFrom(department.DepartmentRoutes)...)
			return err
		}},
		{"position", func() error {
			_, err := web.RegisterHandlers(router, position.GetPositionController(), optionsFrom(position.PositionRoutes)...)
			return err
		}},
		{"member", func() error {
			_, err := web.RegisterHandlers(router, member.GetMemberController(), optionsFrom(member.MemberRoutes)...)
			return err
		}},
		{"invitation", func() error {
			_, err := web.RegisterHandlers(router, invitation.GetInvitationController(), optionsFrom(invitation.InvitationRoutes)...)
			return err
		}},
		{"org_role", func() error {
			_, err := web.RegisterHandlers(router, org_role.GetOrgRoleController(), optionsFrom(org_role.OrgRoleRoutes)...)
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

// optionsFrom 把既有的 di.Routes 路由表转换成框架的注册选项。
//
// 本项目的路由定义历史上保存在 api/<模块>/Routes 变量中（SDK 与文档曾依赖它），
// 而框架以 meta.Meta 标签或注册选项为准。这里复用既有路由表作为运行时覆盖，
// 既不必给数十个请求结构体补标签，也避免维护「两份路由定义」造成漂移。
func optionsFrom(routes di.Routes) []web.RegisterOption {
	opts := make([]web.RegisterOption, 0, len(routes)*2)
	for methodName, route := range routes {
		opts = append(opts, web.WithMethodPath(methodName, route.Path))
		opts = append(opts, web.WithMethodHTTPMethod(methodName, route.Method))
	}
	return opts
}
