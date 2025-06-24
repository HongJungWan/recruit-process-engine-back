package model

import "time"

type Approval struct {
    ApprovalID  int        `db:"approval_id"`
    OfferID     int        `db:"offer_id"`
    ApproverID  int        `db:"approver_id"`
    Status      string     `db:"status"`
    Comment     *string    `db:"comment"`
    RequestedAt time.Time  `db:"requested_at"`
    DecidedAt   *time.Time `db:"decided_at"`
    CreatedBy   string     `db:"created_by"`
}

func (Approval) TableName() string {
    return "approval"
}
