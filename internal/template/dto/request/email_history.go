package request

type CreateEmailHistory struct {
	TemplateID  int    `json:"template_id"    binding:"required"`
	UserID      *int   `json:"user_id,omitempty"`
	ApplicantID *int   `json:"applicant_id,omitempty"`
	OfferID     *int   `json:"offer_id,omitempty"`
	Title       string `json:"title"          binding:"required"`
	Body        string `json:"body"           binding:"required"`
}

type ListEmailHistory struct {
	ApplicantID *int `form:"applicant_id"`
	OfferID     *int `form:"offer_id"`
	Page        int  `form:"page,default=1"`
	Size        int  `form:"size,default=15"`
}
