package model

import "time"

type EmailTemplate struct {
    ID        int            `db:"id"`
    Name      string         `db:"name"`
    Config    []byte         `db:"config"`    // JSONB.RawMessage
    CreatedAt time.Time      `db:"created_at"`
}

func (EmailTemplate) TableName() string {
    return "email_template"
}
