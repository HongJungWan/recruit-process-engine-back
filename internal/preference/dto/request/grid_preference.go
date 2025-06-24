package request

type CreateGridPreference struct {
	GridName string                 `json:"grid_name" binding:"required"`
	Config   map[string]interface{} `json:"config"    binding:"required"`
}

type UpdateGridPreference struct {
	Config map[string]interface{} `json:"config" binding:"required"`
}
