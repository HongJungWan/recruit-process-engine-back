package response

import (
	// 표준 라이브러리
	"time"
)

type ApplicantItem struct {
    ApplicationID int    `json:"application_id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    CurrentStage  string `json:"current_stage"`
}

type ListApplicantsResponse struct {
    Total int              `json:"total"`
    Page  int              `json:"page"`
    Size  int              `json:"size"`
    Items []ApplicantItem  `json:"items"`
}

type ApplicantDetail struct {
    ApplicationID int        `json:"application_id"`
    Name          string     `json:"name"`
    Email         string     `json:"email"`
    Phone         *string    `json:"phone"`
    Education     *string    `json:"education"`
    Experience    *string    `json:"experience"`
    TechStack     *string    `json:"tech_stack"`
    CurrentStage  string     `json:"current_stage"`
    CreatedAt     time.Time  `json:"created_at"`
    CreatedBy     string     `json:"created_by"`
    UpdatedAt     *time.Time `json:"updated_at,omitempty"`
    UpdatedBy     *string    `json:"updated_by,omitempty"`
}

type UpdateStageResponse struct {
    ApplicationID int       `json:"application_id"`
    OldStage      string    `json:"old_stage"`
    NewStage      string    `json:"new_stage"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type BulkUpdateResponse struct {
    Updated int `json:"updated"`
}

type StageHistoryItem struct {
    HistoryID     int        `json:"history_id"`
    Stage         string     `json:"stage"`
    Status        string     `json:"status"`
    CreatedAt     time.Time  `json:"created_at"`
    CreatedBy     string     `json:"created_by"`
    UpdatedAt     *time.Time `json:"updated_at,omitempty"`
    UpdatedBy     *string    `json:"updated_by,omitempty"`
}
