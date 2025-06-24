package response

import "time"

type EmailTemplateItem struct {
    TemplateID int                    `json:"template_id"`
    Name       string                 `json:"name"`
    Config     map[string]interface{} `json:"config"`
    CreatedAt  time.Time              `json:"created_at"`
}

type EmailTemplateDetail = EmailTemplateItem
