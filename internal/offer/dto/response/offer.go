package response

import "time"

type OfferItem struct {
    OfferID       int       `json:"offer_id"`
    ApplicationID int       `json:"application_id"`
    Position      string    `json:"position"`
    Salary        int       `json:"salary"`
    StartDate     time.Time `json:"start_date"`
    Location      string    `json:"location"`
    Status        string    `json:"status"`
}

type ListOffersResponse struct {
    Total int         `json:"total"`
    Items []OfferItem `json:"items"`
}

type ApproverStatus struct {
    ApproverID int     `json:"approver_id"`
    Status     string  `json:"status"`
    Comment    *string `json:"comment,omitempty"`
}

type OfferDetail struct {
    OfferID   int               `json:"offer_id"`
    Approvers []ApproverStatus  `json:"approvers"`
    Status    string            `json:"status"`
}

type SendOfferEmailResponse struct {
    EmailHistoryID int       `json:"email_history_id"`
    SentAt         time.Time `json:"sent_at"`
}