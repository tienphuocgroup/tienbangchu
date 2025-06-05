package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
}

type ServerConfig struct {
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

type LogConfig struct {
	Level string `json:"level"`
}

func Load() *Config {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		Log: LogConfig{
			Level: "info",
		},
	}
}
