package db

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	framecache "github.com/LandcLi/landc-go/frame/pkg/cache"
	framedb "github.com/LandcLi/landc-go/frame/pkg/db"

	"github.com/LandcLi/LandcOrg/utility/config"
)

// 本包是框架基础设施的薄委托层。
//
// 连接建立（driver / dsn / 连接池）统一由 landc-go 框架按 database 段完成，
// 支持 SQLite / MySQL / Postgres；Org 不再自建第二套 gorm 实例，避免
// 「框架连一个库、业务查另一个库」的双连接问题。

// InitDatabase 保留兼容签名：实际连接已在框架 bootstrap 阶段建立。
//
// Deprecated: 请改用 framedb.InitGlobalDBWithDefault()，或直接用装配入口
// （InitDefault / InitWithConfig）。保留此函数仅为兼容既有调用点。
func InitDatabase(_ *config.DatabaseConfig) {}

// InitRedis 保留兼容签名：缓存由框架统一初始化（Redis 不可用时降级为进程内 LRU）。
//
// Deprecated: 请改用 framecache.InitGlobalCacheWithDefault()。
func InitRedis(_ *config.RedisConfig) {}

// GetDB 返回框架全局 gorm 实例。
func GetDB() *gorm.DB { return framedb.GetDB() }

// GetRedis 返回框架的 Redis 客户端；未配置 Redis 或不可用时返回 nil。
//
// 注意：必须在框架初始化后调用。DAO 层应每次调用本函数，而不是在 init() 中
// 捕获结果，否则会永久持有初始化前的 nil（缓存静默失效）。
func GetRedis() *redis.Client { return framecache.GetRedis() }
