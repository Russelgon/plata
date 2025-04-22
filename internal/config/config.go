package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type PostgresConfig struct {
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxIdleTime string `yaml:"connMaxIdleTime"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	HostPrimary     string `yaml:"hostPrimary"`
	HostReplica     string `yaml:"hostReplica"`
	Port            int    `yaml:"port"`
	SSLMode         string `yaml:"sslmode"`
}

type ExchangeConfig struct {
	Timeout time.Duration `yaml:"timeout"`
	APIKey  string        `yaml:"api_key"`
	URL     string        `yaml:"url"`
}

type CronConfig struct {
	Schedule string `yaml:"schedule"`
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Exchange ExchangeConfig `yaml:"exchange"`
	Cron     CronConfig     `yaml:"cron"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	path := "config/config.yaml"
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	defer file.Close()
	return &cfg, nil
}
