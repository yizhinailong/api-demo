package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Log      LogConfig      `toml:"log"`
	Database DatabaseConfig `toml:"database"`
}

type ServerConfig struct {
	Port           string `toml:"port"`
	ReadTimeout    string `toml:"read_timeout"`
	WriteTimeout   string `toml:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes"`
}

type LogConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
	Output string `toml:"output"`
}

type DatabaseConfig struct {
	MySQL    MySQLConfig    `toml:"mysql"`
	Postgres PostgresConfig `toml:"postgres"`
}

type MySQLConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type PostgresConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

var (
	Cfg  *Config
	once sync.Once
	Err  error
)

func GetConfig() *Config {
	once.Do(func() {
		Cfg, Err = Load()
		if Err != nil {
			slog.Error("Failed to load config", "error", Err)
		}
	})

	return Cfg
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Initialize Viper
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("config/")

	// Set default values
	setDefaults(v)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Info("Config file config.toml not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal configuration into struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.max_header_bytes", 1048576)

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")

	// Database defaults
	v.SetDefault("database.mysql.host", "localhost")
	v.SetDefault("database.mysql.port", 3306)
	v.SetDefault("database.mysql.name", "demo")
	v.SetDefault("database.mysql.user", "root")
	v.SetDefault("database.mysql.password", "mysql")

	v.SetDefault("database.postgres.host", "localhost")
	v.SetDefault("database.postgres.port", 5432)
	v.SetDefault("database.postgres.name", "demo")
	v.SetDefault("database.postgres.user", "postgres")
	v.SetDefault("database.postgres.password", "postgresql")
}
