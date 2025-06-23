package config

import (
    "fmt"
    "github.com/joho/godotenv"
    "os"
)

type AppConfig struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    HTTPPort   string
}

var Cfg *AppConfig

func InitConfig() error {
    _ = godotenv.Load(".env")

    cfg := &AppConfig{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "myappdb"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        HTTPPort:   getEnv("HTTP_PORT", "8080"),
    }

    Cfg = cfg
    fmt.Printf("[Config] Loaded: %#v\n", cfg)
    return nil
}

func getEnv(key, defaultVal string) string {
    if v, exists := os.LookupEnv(key); exists {
        return v
    }
    return defaultVal
}
