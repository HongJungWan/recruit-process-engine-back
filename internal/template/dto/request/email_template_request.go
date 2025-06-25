package request

type CreateEmailTemplate struct {
	Name   string                 `json:"name"   binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

type UpdateEmailTemplate struct {
	Name   *string                 `json:"name,omitempty"`
	Config *map[string]interface{} `json:"config,omitempty"`
}
