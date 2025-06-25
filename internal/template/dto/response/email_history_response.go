package response

import (
	// 표준 라이브러리
	"time"
)

type CreateEmailHistoryResponse struct {
    EmailID  int       `json:"email_id"`
    SentAt   time.Time `json:"sent_at"`
}

type EmailHistoryItem struct {
    EmailID int       `json:"email_id"`
    Title   string    `json:"title"`
    SentAt  time.Time `json:"sent_at"`
}

type EmailHistoryDetail struct {
    EmailID       int       `json:"email_id"`
    Title         string    `json:"title"`
    Body          string    `json:"body"`
    SentAt        time.Time `json:"sent_at"`
}
