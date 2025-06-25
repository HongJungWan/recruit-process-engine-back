package response

import (
	// 표준 라이브러리
	"time"
)

type GridPreference struct {
    PreferenceID int                    `json:"preference_id"`        // 설정 ID
    UserID       int                    `json:"user_id"`              // 유저 ID
    GridName     string                 `json:"grid_name"`            // 그리드명
    Config       map[string]interface{} `json:"config"`               // 설정 데이터
    CreatedAt    time.Time              `json:"created_at"`           // 생성일
    CreatedBy    string                 `json:"created_by"`           // 수정자
    UpdatedAt    *time.Time             `json:"updated_at,omitempty"` // 수정일시 (NULL 허용)
    UpdatedBy    *string                `json:"updated_by,omitempty"` // 수정자   (NULL 허용)
}
