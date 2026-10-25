package landcorg

import (
	"fmt"

	"github.com/gin-gonic/gin"

	apiMiddleware "github.com/LandcLi/landc-go/api/middleware/gin"
	framecache "github.com/LandcLi/landc-go/frame/pkg/cache"
	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	framedb "github.com/LandcLi/landc-go/frame/pkg/db"
	"github.com/LandcLi/landc-go/frame/pkg/web"
	ginLogger "github.com/LandcLi/landc-go/log/adapter/gin"

	"github.com/LandcLi/LandcOrg/api"
	_ "github.com/LandcLi/LandcOrg/internal" // 触发默认实现注册（DI）
	"github.com/LandcLi/LandcOrg/model"
	"github.com/LandcLi/LandcOrg/utility/config"
	"github.com/LandcLi/LandcOrg/utility/logger"
)

// InitDefault 从配置文件路径初始化，返回 web.Server。
func InitDefault(configPath string) (*web.Server, error) {
	if err := frameconfig.InitGlobalConfigWithPath(configPath); err != nil {
		return nil, fmt.Errorf("加载框架配置失败: %w", err)
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return InitWithConfig(cfg)
}

// InitWithConfig 从配置对象初始化，返回 web.Server。
// 前置条件：调用方已完成框架配置加载（InitDefault 或宿主 bootstrap）。
func InitWithConfig(cfg *config.Config) (*web.Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置不能为空")
	}
	if cfg.Server.GinMode != "" {
		gin.SetMode(cfg.Server.GinMode)
	}

	server := web.New()
	engine := server.Engine()

	if err := initInfrastructure(cfg); err != nil {
		return nil, err
	}

	// api 层中间件（trace / 请求日志等）
	engine.Use(apiMiddleware.GinMiddleware())
	ginLogger.UseWithGin(engine, logger.GetLogger())

	if cfg.Server.UseDefaultRoute {
		if err := api.RegisterDefaultRoutes(server); err != nil {
			return nil, err
		}
	}
	return server, nil
}

// InitWithConfigAndGin 使用已有 Gin 引擎初始化（库模式 / 嵌入到宿主）。
// 不挂载全局中间件：中间件链由宿主决定。
func InitWithConfigAndGin(cfg *config.Config, ginEngine *gin.Engine) error {
	if cfg == nil {
		return fmt.Errorf("配置不能为空")
	}
	if ginEngine == nil {
		return fmt.Errorf("gin 引擎不能为空")
	}
	if cfg.Server.GinMode != "" {
		gin.SetMode(cfg.Server.GinMode)
	}
	if err := initInfrastructure(cfg); err != nil {
		return err
	}
	if cfg.Server.UseDefaultRoute {
		return api.RegisterHandlers(ginEngine)
	}
	return nil
}

// Models 返回本系统负责的持久化模型（模型清单唯一来源）。
func Models() []interface{} {
	return []interface{}{
		&model.Department{},
		&model.Position{},
		&model.OrgMember{},
		&model.OrgInvitation{},
		&model.DepartmentRole{},
		&model.PositionRole{},
	}
}

// EnsureSchema 按配置决定是否自动建表（受 database.auto_migrate 控制），幂等。
func EnsureSchema() error {
	if err := framedb.AutoMigrateIfEnabled(Models()...); err != nil {
		return fmt.Errorf("自动建表失败: %w", err)
	}
	return nil
}

// initInfrastructure 初始化日志 / DB / 缓存 / 表结构。
// DB 与缓存统一由框架按 database / redis 段建立（SQLite、MySQL、Postgres 均可）。
func initInfrastructure(cfg *config.Config) error {
	logger.InitLogger(config.FacadeLogConfig())

	if err := framedb.InitGlobalDBWithDefault(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}
	if err := framecache.InitGlobalCacheWithDefault(); err != nil {
		return fmt.Errorf("初始化缓存失败: %w", err)
	}
	return EnsureSchema()
}
