package model

import (
	// 표준 라이브러리
	"time"
)

type User struct {
    UserID    int        `db:"user_id"    json:"user_id"`             // 유저 ID
    LoginID   string     `db:"login_id"   json:"login_id"`            // 로그인 아이디
    Password  string     `db:"login_pw"   json:"-"`                   // 로그인 패스워드
    Name      string     `db:"name"       json:"name"`                // 유저 이름
    Email     string     `db:"email"      json:"email"`               // 유저 이메일 주소
    Role      string     `db:"role"       json:"role"`                // 유저 역할
    CreatedAt time.Time  `db:"created_at" json:"created_at"`          // 생성일시
    CreatedBy string     `db:"created_by" json:"created_by"`          // 생성자
    UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`// 수정일시 (NULL 허용)
    UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`// 수정자   (NULL 허용)
}

func (User) TableName() string {
    return "users"
}
