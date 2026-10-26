// Package middleware 提供 LandcOrg 的 HTTP 中间件。
package middleware

import (
	"crypto/subtle"

	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	framemiddleware "github.com/LandcLi/landc-go/frame/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// ServiceKeyHeader 服务间调用凭证的请求头名。
const ServiceKeyHeader = "X-Service-Key"

// ServiceKeyConfigPath 服务密钥在配置文件中的位置。
//
// 放在 other.service_auth.key：框架的 Config 结构没有这个字段，
// 通过 GetOtherString 读取 other 段的自由扩展配置。
var ServiceKeyConfigPath = []string{"service_auth", "key"}

// Auth 是本系统的鉴权中间件。
//
// 历史缺陷：InitWithConfig 此前只挂了 trace 与请求日志中间件，**完全没有鉴权** ——
// 任何人只要知道地址就能创建部门、添加成员、分配组织角色。同为内网服务的
// LandcPermission 是挂了框架 middleware.Auth() 的，本系统属于遗漏。
//
// 现在的策略：
//   - 携带正确 X-Service-Key 的内部服务调用直接放行（Permission 需要按 user_id
//     查用户继承角色，那是没有用户 token 的服务间调用）；
//   - 其余请求一律走框架 JWT 校验，健康检查等探针端点由框架默认跳过。
//
// 未配置服务密钥时恒不放行服务间调用，即「漏配只会更严格」。
func Auth() gin.HandlerFunc {
	jwtAuth := framemiddleware.Auth()

	return func(c *gin.Context) {
		if isServiceCall(c) {
			c.Next()
			return
		}
		jwtAuth(c)
	}
}

func isServiceCall(c *gin.Context) bool {
	key := serviceKey()
	if key == "" {
		return false
	}
	got := c.GetHeader(ServiceKeyHeader)
	if got == "" {
		return false
	}
	// 恒定时间比较，避免通过响应耗时逐字节猜密钥
	return subtle.ConstantTimeCompare([]byte(got), []byte(key)) == 1
}

// serviceKey 读取服务密钥。每次请求读取而非缓存，便于配置热加载与测试替换。
func serviceKey() string {
	fc := frameconfig.GetConfig()
	if fc == nil {
		return ""
	}
	return fc.GetOtherString(ServiceKeyConfigPath...)
}
