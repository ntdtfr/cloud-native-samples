// config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port         string
		ReadTimeout  int
		WriteTimeout int
	}

	MongoDB struct {
		URI      string
		Database string
	}

	Redis struct {
		Address  string
		Password string
		DB       int
	}

	RabbitMQ struct {
		URI      string
		Exchange string
	}

	Auth struct {
		JWTSecret     string
		TokenDuration int
	}

	RateLimit struct {
		Requests int
		Duration int
	}

	LogLevel string
}

func Load() (*Config, error) {
	var config Config

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.readTimeout", 10)
	viper.SetDefault("server.writeTimeout", 10)
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "product_service")
	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("rabbitmq.uri", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("rabbitmq.exchange", "product_exchange")
	viper.SetDefault("auth.jwtSecret", "secret")
	viper.SetDefault("auth.tokenDuration", 3600)
	viper.SetDefault("rateLimit.requests", 100)
	viper.SetDefault("rateLimit.duration", 60)
	viper.SetDefault("logLevel", "info")

	// Load from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Override with environment variables
	viper.AutomaticEnv()

	// Map environment variables
	overrideWithEnv("SERVER_PORT", "server.port")
	overrideWithEnv("SERVER_READ_TIMEOUT", "server.readTimeout")
	overrideWithEnv("SERVER_WRITE_TIMEOUT", "server.writeTimeout")
	overrideWithEnv("MONGODB_URI", "mongodb.uri")
	overrideWithEnv("MONGODB_DATABASE", "mongodb.database")
	overrideWithEnv("REDIS_ADDRESS", "redis.address")
	overrideWithEnv("REDIS_PASSWORD", "redis.password")
	overrideWithEnvInt("REDIS_DB", "redis.db")
	overrideWithEnv("RABBITMQ_URI", "rabbitmq.uri")
	overrideWithEnv("RABBITMQ_EXCHANGE", "rabbitmq.exchange")
	overrideWithEnv("AUTH_JWT_SECRET", "auth.jwtSecret")
	overrideWithEnvInt("AUTH_TOKEN_DURATION", "auth.tokenDuration")
	overrideWithEnvInt("RATE_LIMIT_REQUESTS", "rateLimit.requests")
	overrideWithEnvInt("RATE_LIMIT_DURATION", "rateLimit.duration")
	overrideWithEnv("LOG_LEVEL", "logLevel")

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

func overrideWithEnv(envKey, configKey string) {
	if val, exists := os.LookupEnv(envKey); exists {
		viper.Set(configKey, val)
	}
}

func overrideWithEnvInt(envKey, configKey string) {
	if val, exists := os.LookupEnv(envKey); exists {
		if intVal, err := strconv.Atoi(val); err == nil {
			viper.Set(configKey, intVal)
		}
	}
}
