package model

import (
	"encoding/json"
	"time"
)

type GridPreference struct {
    PreferenceID int             `db:"preference_id"`
    UserID       int             `db:"user_id"`
    GridName     string          `db:"grid_name"`
    Config       json.RawMessage `db:"config"`       // JSONB
    CreatedAt    time.Time       `db:"created_at"`
    CreatedBy    string          `db:"created_by"`
    UpdatedAt    *time.Time      `db:"updated_at"`
    UpdatedBy    *string         `db:"updated_by"`
}

func (GridPreference) TableName() string {
    return "user_grid_preference"
}
