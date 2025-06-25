package model

import (
	// 표준 라이브러리
	"time"
)

type EmailHistory struct {
    EmailID       int        `db:"email_id"`       // 이력 ID
    UserID        *int       `db:"user_id"`        // 발송자 ID
    ApplicationID *int       `db:"application_id"` // 지원자 ID
    OfferID       *int       `db:"offer_id"`       // 오퍼 ID
    TemplateID    int        `db:"template_id"`    // 템플릿 ID
    Title         string     `db:"title"`          // 제목
    Body          string     `db:"body"`           // 본문
    CreatedAt     time.Time  `db:"created_at"`     // 발송일
    CreatedBy     string     `db:"created_by"`     // 발송자
}

func (EmailHistory) TableName() string {
    return "email_history"
}
