package api

import (
	"github.com/LandcLi/LandcOrg/api/department"
	"github.com/LandcLi/LandcOrg/api/invitation"
	"github.com/LandcLi/LandcOrg/api/member"
	"github.com/LandcLi/LandcOrg/api/org_role"
	"github.com/LandcLi/LandcOrg/api/position"
	"github.com/gin-gonic/gin"
)

// RegisterDefaultRoutes registers all default routes using Gateway auto-generation.
func RegisterDefaultRoutes(r *gin.Engine) {
	// 健康检查已由框架自动注册（config.use_default_routes: true）
	department.DepartmentGateway.RegisterRoutes(r)
	position.PositionGateway.RegisterRoutes(r)
	member.MemberGateway.RegisterRoutes(r)
	invitation.InvitationGateway.RegisterRoutes(r)
	org_role.OrgRoleGateway.RegisterRoutes(r)
}
