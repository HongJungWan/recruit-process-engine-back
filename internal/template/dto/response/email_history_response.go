package response

import (
	// 표준 라이브러리
	"time"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
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

func ToEmailHistoryDetail(h *model.EmailHistory) EmailHistoryDetail {
    return EmailHistoryDetail{
        EmailID: h.EmailID,
        Title:   h.Title,
        Body:    h.Body,
        SentAt:  h.CreatedAt,
    }
}