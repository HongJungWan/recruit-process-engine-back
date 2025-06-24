package response

import "time"

type GridPreference struct {
    PreferenceID int                    `json:"preference_id"`
    UserID       int                    `json:"user_id"`
    GridName     string                 `json:"grid_name"`
    Config       map[string]interface{} `json:"config"`
    CreatedAt    time.Time              `json:"created_at"`
    CreatedBy    string                 `json:"created_by"`
    UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
    UpdatedBy    *string                `json:"updated_by,omitempty"`
}
