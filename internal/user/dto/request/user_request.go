package request

type Credentials struct {
	LoginId  string `json:"login_id" binding:"required"` // 사용자 이름
	Password string `json:"password" binding:"required"` // 비밀번호
}
