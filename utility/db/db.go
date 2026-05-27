package db

import (
	"github.com/LandcLi/LandcOrg/utility/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	database    *gorm.DB
	redisClient *redis.Client
)

func InitDatabase(cfg *config.DatabaseConfig) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	database = db
}

func GetDB() *gorm.DB { return database }

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: cfg.Addr, Password: cfg.Password, DB: cfg.DB,
	})
	return redisClient, nil
}

func GetRedis() *redis.Client { return redisClient }
