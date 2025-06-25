package model

import (
	// 표준 라이브러리
	"encoding/json"
	"time"
)

type GridPreference struct {
    PreferenceID int             `db:"preference_id"` // 설정 ID
    UserID       int             `db:"user_id"`       // 유저 ID
    GridName     string          `db:"grid_name"`     // 그리드명
    Config       json.RawMessage `db:"config"`        // 설정 데이터, JSONB
    CreatedAt    time.Time       `db:"created_at"`    // 생성일
    CreatedBy    string          `db:"created_by"`    // 생성자
    UpdatedAt    *time.Time      `db:"updated_at"`    // 수정일 (NULL 허용)
    UpdatedBy    *string         `db:"updated_by"`    // 수정자 (NULL 허용)
}

func (GridPreference) TableName() string {
    return "user_grid_preference"
}
