package session

import (
	"database/sql"
	"time"

	ginAdapter "github.com/39george/scs_gin_adapter"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

var (
    Manager *scs.SessionManager // 세션 매니저
    Adapter *ginAdapter.GinAdapter // 세션 매니저 어댑터,
)

// 세션 매니저 초기화
func InitSession(db *sql.DB) {
    m := scs.New()
    m.Lifetime = 24 * time.Hour
    m.IdleTimeout = 30 * time.Minute
    m.Cookie.HttpOnly = true
    m.Cookie.Secure = false    // HTTPS 환경에선 true 권장
    m.Store = postgresstore.New(db)  // PostgreSQL 스토어 사용

    Manager = m
    Adapter = ginAdapter.New(m)
}
