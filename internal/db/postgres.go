package db

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
    "log"
    "github.com/HongJungWan/recruit-process-engine-back/internal/config"
)

var DB *sqlx.DB

func InitDB() error {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        config.Cfg.DBUser,
        config.Cfg.DBPassword,
        config.Cfg.DBHost,
        config.Cfg.DBPort,
        config.Cfg.DBName,
        config.Cfg.DBSSLMode,
    )

    db, err := sqlx.Open("pgx", dsn)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)

    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    DB = db
    log.Println("[DB] Connected to PostgreSQL successfully")
    return nil
}
