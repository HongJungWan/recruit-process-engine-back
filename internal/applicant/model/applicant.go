package model

import (
	// 표준 라이브러리
	"time"
)

type Applicant struct {
    ApplicationID int        `db:"application_id"` // 지원자 ID
    Name          string     `db:"name"`           // 지원자 이름
    Email         string     `db:"email"`          // 지원자 이메일
    Phone         *string    `db:"phone"`          // 전화번호
    Education     *string    `db:"education"`      // 학력
    Experience    *string    `db:"experience"`     // 경력
    TechStack     *string    `db:"tech_stack"`     // 기술 스택
    CurrentStage  string     `db:"current_stage"`  // 전형 단계
    CreatedAt     time.Time  `db:"created_at"`     // 생성일
    CreatedBy     string     `db:"created_by"`     // 생성자
    UpdatedAt     *time.Time `db:"updated_at"`     // 수정일 (NULL 허용)
    UpdatedBy     *string    `db:"updated_by"`     // 수정자 (NULL 허용)
}

func (Applicant) TableName() string {
    return "applicant"
}
