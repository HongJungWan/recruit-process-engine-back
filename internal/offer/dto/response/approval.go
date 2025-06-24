package response

import "time"

type ApprovalStatus struct {
    ApprovalID  int        `json:"approval_id"`
    ApproverID  int        `json:"approver_id"`
    Status      string     `json:"status"`
    Comment     *string    `json:"comment,omitempty"`
    RequestedAt time.Time  `json:"requested_at"`
    DecidedAt   *time.Time `json:"decided_at,omitempty"`
}

type ApprovalHistoryItem struct {
    ApprovalID  int        `json:"approval_id"`
    ApproverID  int        `json:"approver_id"`
    Status      string     `json:"status"`
    Comment     *string    `json:"comment,omitempty"`
    RequestedAt time.Time  `json:"requested_at"`
    DecidedAt   *time.Time `json:"decided_at,omitempty"`
}

type CreateApprovalsResponse []ApprovalHistoryItem

type ProcessApprovalResponse struct {
    ApprovalID int       `json:"approval_id"`
    Status     string    `json:"status"`
    DecidedAt  time.Time `json:"decided_at"`
}

type SendOfferResponse struct {
    EmailHistoryID int       `json:"email_history_id"`
    SentAt         time.Time `json:"sent_at"`
}
