package model

import "time"

type EmailHistory struct {
    EmailID       int        `db:"email_id"`
    UserID        *int       `db:"user_id"`
    ApplicationID *int       `db:"application_id"`
    OfferID       *int       `db:"offer_id"`
    TemplateID    int        `db:"template_id"`
    Title         string     `db:"title"`
    Body          string     `db:"body"`
    CreatedAt     time.Time  `db:"created_at"`
    CreatedBy     string     `db:"created_by"`
}

func (EmailHistory) TableName() string {
    return "email_history"
}
