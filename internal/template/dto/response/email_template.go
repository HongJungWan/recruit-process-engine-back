package response

import (
	// 표준 라이브러리
	"time"
)

type EmailTemplateItem struct {
    TemplateID int                    `json:"template_id"`
    Name       string                 `json:"name"`
    Config     map[string]interface{} `json:"config"`
    CreatedAt  time.Time              `json:"created_at"`
}

type EmailTemplateDetail = EmailTemplateItem
