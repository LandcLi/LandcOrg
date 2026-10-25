package config

import (
	"os"
	"time"

	frameconfig "github.com/LandcLi/landc-go/frame/pkg/config"
	"github.com/LandcLi/landc-go/log/facade"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	// 注意：log 段由框架（landc-go/frame/pkg/config）解析。
	// 本结构不再重复定义，避免同一个配置项出现不兼容的类型
	// （框架为字符串级别，facade 为 int 枚举，重复解析必然一方失败）。
}

type ServerConfig struct {
	Port            int    `yaml:"port"`
	GinMode         string `yaml:"gin_mode"`
	EnableCORS      bool   `yaml:"enable_cors"`
	UseDefaultRoute bool   `yaml:"use_default_route"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
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
