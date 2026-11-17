package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	NSQ      NSQConfig      `yaml:"nsq"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type NSQConfig struct {
	NSQDTCPAddr string `yaml:"nsqd_tcp_addr"`
	Topic       string `yaml:"topic"`
	Channel     string `yaml:"channel"`
	MaxAttempts uint16 `yaml:"max_attempts"`
	Concurrency int    `yaml:"concurrency"`
}

type RedisConfig struct {
	Addr                  string `yaml:"addr"`
	LockTTLSeconds        int    `yaml:"lock_ttl_seconds"`
	IdempotencyTTLSeconds int    `yaml:"idempotency_ttl_seconds"`
}

// LoadByEnv loads files/etc/<app>/<app>.<env>.yaml using APP_ENV (defaults to "local").
func LoadByEnv(appName string) (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	path := fmt.Sprintf("files/etc/%s/%s.%s.yaml", appName, appName, env)
	return Load(path)
}

func Load(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
