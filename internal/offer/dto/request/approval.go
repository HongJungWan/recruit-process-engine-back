package request

type CreateApprovalsRequest struct {
	ApproverIDs []int `json:"approver_ids" binding:"required"`
}

type ProcessApprovalRequest struct {
	Status  string  `json:"status"  binding:"required"`
	Comment *string `json:"comment,omitempty"`
}
