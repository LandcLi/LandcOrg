package api

import (
	"github.com/LandcLi/LandcOrg/api/department"
	"github.com/LandcLi/LandcOrg/api/member"
	"github.com/LandcLi/LandcOrg/api/position"
	"github.com/gin-gonic/gin"
)

// RegisterDefaultRoutes registers all default routes using Gateway auto-generation.
func RegisterDefaultRoutes(r *gin.Engine) {
	r.GET("/health", HealthCheck)

	department.DepartmentGateway.RegisterRoutes(r)
	position.PositionGateway.RegisterRoutes(r)
	member.MemberGateway.RegisterRoutes(r)
}

// HealthCheck health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
