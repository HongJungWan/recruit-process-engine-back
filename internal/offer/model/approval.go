package model

import (
	// 표준 라이브러리
	"time"
)

type Approval struct {
    ApprovalID  int        `db:"approval_id"`  // 결재 ID
    OfferID     int        `db:"offer_id"`     // 오퍼 ID
    ApproverID  int        `db:"approver_id"`  // 결재자 ID
    Status      string     `db:"status"`       // 결재 상태
    Comment     *string    `db:"comment"`      // 결재 코멘트
    RequestedAt time.Time  `db:"requested_at"` // 요청 시각
    DecidedAt   *time.Time `db:"decided_at"`   // 결재 시각
    CreatedAt     time.Time  `db:"created_at"` // 생성일
    CreatedBy   string     `db:"created_by"`   // 생성자
}

func (Approval) TableName() string {
    return "approval"
}
