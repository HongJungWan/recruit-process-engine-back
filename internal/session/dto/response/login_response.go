package response

import "time"

type UserInfo struct {
    UserID    int       `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    CreatedBy string    `json:"created_by"`
}

type LoginResponse struct {
    SessionToken string   `json:"session_token"`
    User         UserInfo `json:"user"`
}