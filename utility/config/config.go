package config

import (
	"os"

	"github.com/LandcLi/landc-go/log/facade"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig     `yaml:"server"`
	Database DatabaseConfig   `yaml:"database"`
	Redis    RedisConfig      `yaml:"redis"`
	Log      facade.LogConfig `yaml:"log"`
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
