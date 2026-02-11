package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Logger      LoggerConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MaxIdle  int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string
	OutputPath string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "testmesh")
	viper.SetDefault("database.password", "testmesh")
	viper.SetDefault("database.dbname", "testmesh")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_conns", 25)
	viper.SetDefault("database.max_idle", 5)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.output_path", "stdout")

	// Auto-load environment variables
	viper.AutomaticEnv()

	// Try to load config file (optional)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.ReadInConfig() // Ignore error if config file not found

	readTimeout, _ := time.ParseDuration(viper.GetString("server.read_timeout"))
	writeTimeout, _ := time.ParseDuration(viper.GetString("server.write_timeout"))

	cfg := &Config{
		Environment: viper.GetString("environment"),
		Server: ServerConfig{
			Port:         viper.GetInt("server.port"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			SSLMode:  viper.GetString("database.sslmode"),
			MaxConns: viper.GetInt("database.max_conns"),
			MaxIdle:  viper.GetInt("database.max_idle"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Logger: LoggerConfig{
			Level:      viper.GetString("logger.level"),
			OutputPath: viper.GetString("logger.output_path"),
		},
	}

	return cfg, nil
}
