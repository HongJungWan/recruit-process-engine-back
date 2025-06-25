package model

import (
	// 표준 라이브러리
	"time"
)

type Offer struct {
    OfferID       int        `db:"offer_id"`
    UserID        int        `db:"user_id"`
    ApplicationID int        `db:"application_id"`
    Position      string     `db:"position"`
    Salary        int        `db:"salary"`
    StartDate     time.Time  `db:"start_date"`
    Location      string     `db:"location"`
    Benefits      string     `db:"benefits"`
    LetterContent string     `db:"letter_content"`
    Status        string     `db:"status"`
    ApprovedAt    *time.Time `db:"approved_at"`
    SentAt        time.Time  `db:"sent_at"`
    CreatedAt     time.Time  `db:"created_at"`
    CreatedBy     string     `db:"created_by"`
    UpdatedAt     *time.Time `db:"updated_at"`
    UpdatedBy     *string    `db:"updated_by"`
}

func (Offer) TableName() string {
    return "offer"
}
