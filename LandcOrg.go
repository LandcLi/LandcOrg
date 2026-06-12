package landcorg

import (
	"fmt"

	"github.com/LandcLi/LandcOrg/api"
	_ "github.com/LandcLi/LandcOrg/internal"
	"github.com/LandcLi/LandcOrg/utility/config"
	"github.com/LandcLi/LandcOrg/utility/db"
	"github.com/LandcLi/LandcOrg/utility/logger"
	apiMiddleware "github.com/LandcLi/landc-go/api/middleware/gin"
	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	"github.com/LandcLi/landc-go/frame/pkg/web"
	ginLogger "github.com/LandcLi/landc-go/log/adapter/gin"
	"github.com/gin-gonic/gin"
)

// InitDefault initializes from config file path.
func InitDefault(configPath string) (*web.Server, error) {
	if err := frameconfig.InitGlobalConfigWithPath(configPath); err != nil {
		return nil, fmt.Errorf("加载框架配置失败: %w", err)
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return initWithConfig(cfg)
}

// InitWithConfig initializes from config object.
func InitWithConfig(cfg *config.Config) (*web.Server, error) {
	return initWithConfig(cfg)
}

func initWithConfig(cfg *config.Config) (*web.Server, error) {
	// 使用框架 web.Server，自动注册健康检查和超时中间件
	server := web.New()
	engine := server.Engine()

	if cfg.Server.GinMode != "" {
		gin.SetMode(cfg.Server.GinMode)
	}

	logger.InitLogger(&cfg.Log)
	db.InitRedis(&cfg.Redis)
	db.InitDatabase(&cfg.Database)

	// 应用 api 层面的中间件
	engine.Use(apiMiddleware.GinMiddleware())
	ginLogger.UseWithGin(engine, logger.GetLogger())

	// 注册业务路由
	if cfg.Server.UseDefaultRoute {
		api.RegisterDefaultRoutes(engine)
	}

	return server, nil
}
