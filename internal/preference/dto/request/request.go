package request

type CreateGridPreference struct {
	GridName string                 `json:"grid_name" binding:"required"` // 그리드명
	Config   map[string]interface{} `json:"config"    binding:"required"` // 설정 데이터
}

type UpdateGridPreference struct {
	Config map[string]interface{} `json:"config" binding:"required"` // 설정 데이터
}
