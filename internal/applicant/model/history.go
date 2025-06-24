package model

import "time"

type StageHistory struct {
    HistoryID     int        `db:"history_id"`
    ApplicationID int        `db:"application_id"`
    Stage         string     `db:"stage"`
    Status        string     `db:"status"`
    CreatedAt     time.Time  `db:"created_at"`
    CreatedBy     string     `db:"created_by"`
    UpdatedAt     *time.Time `db:"updated_at"`
    UpdatedBy     *string    `db:"updated_by"`
}

func (StageHistory) TableName() string {
    return "applicant_stage_history"
}
