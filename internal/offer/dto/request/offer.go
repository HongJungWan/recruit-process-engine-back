package request

import "time"

type CreateOfferRequest struct {
    ApplicationID int       `json:"application_id" binding:"required"`
    Position      string    `json:"position"       binding:"required"`
    Salary        int       `json:"salary"         binding:"required"`
    StartDate     time.Time `json:"start_date"     binding:"required"`
    Location      string    `json:"location"       binding:"required"`
    Benefits      string    `json:"benefits"       binding:"required"`
    LetterContent string    `json:"letter_content" binding:"required"`
}

type ListOffersRequest struct {
    Status string `form:"status"`
    Page   int    `form:"page,default=1"`
    Size   int    `form:"size,default=15"`
}

type SendOfferEmailRequest struct{}