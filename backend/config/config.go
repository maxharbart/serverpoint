package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	AESKey     string
	Port       string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "serverpoint"),
		DBPassword: getEnv("DB_PASSWORD", "serverpoint"),
		DBName:     getEnv("DB_NAME", "serverpoint"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
		AESKey:     getEnv("AES_KEY", "0123456789abcdef0123456789abcdef"),
		Port:       getEnv("PORT", "8080"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
