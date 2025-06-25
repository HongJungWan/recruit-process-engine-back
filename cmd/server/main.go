// @title       Recruit Process Engine API
// @version     1.0
// @description 채용 프로세스 엔진 백엔드 API 문서
// @contact.name 홍정완
// @contact.email test@example.com
// @host        localhost:8080
// @BasePath    /api/v1
package main

import (
	// 표준 라이브러리
	"log"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/configs/config"
	"github.com/HongJungWan/recruit-process-engine-back/configs/db"
	"github.com/HongJungWan/recruit-process-engine-back/internal/router"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
)

func main() {
    // 1. 설정 로드
    if err := config.InitConfig(); err != nil {
        log.Fatalf("설정 로드 실패: %v", err)
    }

    // 2. DB 초기화 (sqlx + PostgreSQL)
    if err := db.InitDB(); err != nil {
        log.Fatalf("DB 연결 실패: %v", err)
    }

    // 3. 세션 매니저 초기화
    sqlDB := db.DB.DB
    session.InitSession(sqlDB)

    // 4. 라우터 초기화
    r := router.InitRouter(db.DB)

    // 5. HTTP 서버 구동
    addr := ":" + config.Cfg.HTTPPort
    log.Printf("[서버] Listening on %s\n", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("서버 실행 실패: %v", err)
    }
}
