package config

import (
	"time"

	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	"github.com/LandcLi/landc-go/log/facade"
)

// Config 是本系统对框架配置的视图。
//
// 解析由框架独占完成（landc-go/frame/pkg/config 是配置格式的唯一权威），
// 本包只做字段映射。
//
// 历史缺陷：这里曾用一套独立的 yaml 结构重复解析同一份配置文件，且字段名
// 与框架不一致 —— 本包 `use_default_route`（单数）vs 框架
// `server.use_default_routes`（复数）。后果是：按框架文档写配置时，
// 本包读到 false，业务控制器整体不注册，所有 /api/* 请求落到 NoRoute
// 返回 200 + 空 data。接口看上去「通了」，实际什么都没执行，
// 排查时极易被 200 误导。改为单一解析源后这类漂移不可能再出现。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

// ServerConfig 只保留本系统真正读取的字段。
// 监听地址与端口由框架 server.addr / server.port 决定（web.New 读取全局配置），
// 本系统不重复持有，避免出现「配了却不生效」的假字段。
type ServerConfig struct {
	// GinMode gin 运行模式（debug / release / test），取自 other.server.gin_mode
	GinMode string
	// UseDefaultRoute 是否注册本系统控制器路由，取自框架 server.use_default_routes
	UseDefaultRoute bool
	// EnableCORS 配置视图。middleware.CORS() 的启用策略由框架决定。
	EnableCORS bool
}

// DatabaseConfig 数据库配置视图（连接由框架 pkg/db 建立）。
type DatabaseConfig struct {
	Driver string
	DSN    string
}

// RedisConfig 缓存配置视图（连接由框架 pkg/cache 建立）。
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// LoadConfig 从配置文件加载配置。
// 解析交给框架，保证与框架读取到的是同一份语义（同样的默认值、同样的字段名）。
func LoadConfig(path string) (*Config, error) {
	fc, err := frameconfig.LoadConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	return FromFramework(fc), nil
}

// FromFramework 把已加载的框架配置转换为本系统视图。
// 用于调用方已完成框架配置加载的场景（InitDefault / 宿主 bootstrap），
// 避免重复读文件。
func FromFramework(fc *frameconfig.Config) *Config {
	if fc == nil {
		return &Config{}
	}
	return &Config{
		Server: ServerConfig{
			GinMode:         fc.GetOtherString("server", "gin_mode"),
			UseDefaultRoute: fc.Server.UseDefaultRoutes,
			EnableCORS:      enableCORS(fc),
		},
		Database: DatabaseConfig{
			Driver: fc.Database.Driver,
			DSN:    fc.Database.DSN,
		},
		Redis: RedisConfig{
			Addr:     fc.Redis.Addr,
			Password: fc.Redis.Password,
			DB:       fc.Redis.DB,
		},
	}
}

// enableCORS 读取 CORS 开关。优先 other.server.enable_cors，兼容顶层写法。
func enableCORS(fc *frameconfig.Config) bool {
	if v, ok := fc.GetOtherBool("server", "enable_cors"); ok {
		return v
	}
	if v, ok := fc.GetOtherBool("enable_cors"); ok {
		return v
	}
	return false
}

// FacadeLogConfig 从框架配置的 log 段构造业务侧日志配置。
// 未加载框架配置时返回框架默认值，保证日志系统始终可用。
func FacadeLogConfig() *facade.LogConfig {
	out := &facade.LogConfig{
		Provider:   "console",
		Name:       "landc-org",
		Level:      facade.InfoLevel,
		Format:     "text",
		OutputPath: "stdout",
	}
	fc := frameconfig.GetConfig()
	if fc == nil {
		return out
	}
	out.Level = parseLogLevel(fc.Log.Level)
	if fc.Log.Format != "" {
		out.Format = fc.Log.Format
	}
	if fc.Log.Output != "" {
		out.OutputPath = fc.Log.Output
	}
	if fc.Log.MaxSize > 0 {
		out.MaxLogSize = int64(fc.Log.MaxSize)
	}
	if fc.Log.MaxBackups > 0 {
		out.MaxLogFiles = fc.Log.MaxBackups
	}
	if fc.Log.MaxAge > 0 {
		out.MaxLogAge = time.Duration(fc.Log.MaxAge) * 24 * time.Hour
	}
	return out
}

func parseLogLevel(level string) facade.LogLevel {
	switch level {
	case "debug":
		return facade.DebugLevel
	case "warn", "warning":
		return facade.WarnLevel
	case "error":
		return facade.ErrorLevel
	default:
		return facade.InfoLevel
	}
}
