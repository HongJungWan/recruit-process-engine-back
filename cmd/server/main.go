package main

import (
	"log"

	"github.com/HongJungWan/recruit-process-engine-back/internal/config"
	"github.com/HongJungWan/recruit-process-engine-back/internal/db"
	"github.com/HongJungWan/recruit-process-engine-back/internal/router"
)

func main() {
    // 1) 설정 로드
    if err := config.InitConfig(); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2) DB 초기화 (sqlx + PostgreSQL)
    if err := db.InitDB(); err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }

    // 3) 라우터 초기화
    r := router.InitRouter(db.DB)

    // 4) HTTP 서버 구동
    addr := ":" + config.Cfg.HTTPPort
    log.Printf("[Server] Listening on %s\n", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
