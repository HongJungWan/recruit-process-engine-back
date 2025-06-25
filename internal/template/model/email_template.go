package model

import (
	// 표준 라이브러리
	"time"
)

type EmailTemplate struct {
    ID        int            `db:"id"`         // 템플릿 ID
    Name      string         `db:"name"`       // 식별자
    Config    []byte         `db:"config"`     // 설정 데이터, JSONB
    CreatedAt time.Time      `db:"created_at"` // 생성일
    // TODO: 생성자
}

func (EmailTemplate) TableName() string {
    return "email_template"
}
