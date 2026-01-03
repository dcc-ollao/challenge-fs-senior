package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (Config, error) {
	cfg := Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    os.Getenv("DB_HOST"),
		DBUser:    os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}

	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT: %q", dbPortStr)
	}
	cfg.DBPort = dbPort

	missing := []string{}
	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %v", missing)
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
