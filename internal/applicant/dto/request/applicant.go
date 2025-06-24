package request

type ListApplicantsRequest struct {
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=15"`
	Stage   string `form:"stage"`
	Keyword string `form:"keyword"`
}

type UpdateStageRequest struct {
	Stage  string `json:"stage"  binding:"required"`
	Reason string `json:"reason,omitempty"`
}

type BulkUpdateStageRequest struct {
	IDs    []int  `json:"ids"   binding:"required"`
	Stage  string `json:"stage" binding:"required"`
	Reason string `json:"reason,omitempty"`
}
