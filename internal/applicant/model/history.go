package model

import "time"

type StageHistory struct {
    HistoryID     int        `db:"history_id"`     // 히스토리 ID
    ApplicationID int        `db:"application_id"` // 지원자 ID
    Stage         string     `db:"stage"`          // 단계
    Status        string     `db:"status"`         // 결과
    CreatedAt     time.Time  `db:"created_at"`     // 생성일
    CreatedBy     string     `db:"created_by"`     // 생성자
    UpdatedAt     *time.Time `db:"updated_at"`     // 수정일 (NULL 허용)
    UpdatedBy     *string    `db:"updated_by"`     // 수정자 (NULL 허용)
}

func (StageHistory) TableName() string {
    return "applicant_stage_history"
}
