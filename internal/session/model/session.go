package model

import (
	"encoding/json"
	"time"
)

type Session struct {
    SessionID    string          `db:"session_id"    json:"session_id"`    // 세션 PK (UUID)
    UserID       int             `db:"user_id"       json:"user_id"`       // 유저 ID (FK)
    SessionToken string          `db:"session_token" json:"session_token"` // 클라이언트용 토큰
    Data         json.RawMessage `db:"data"          json:"data"`          // 추가 세션 데이터(JSONB)
    CreatedAt    time.Time       `db:"created_at"    json:"created_at"`    // 생성일시
    ExpiresAt    time.Time       `db:"expires_at"    json:"expires_at"`    // 만료일시
}

func (Session) TableName() string {
    return "sessions"
}
