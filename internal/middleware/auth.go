package middleware

import (
	"net/http"
	"strings"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session/repository"
	"github.com/gin-gonic/gin"
)

// feat: 인증 미들웨어, 로그인(세션 생성)된 사용자만 접근 가능
func Auth(sessionRepo repository.SessionRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")

        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
            c.Abort()
            return
        }

        sess, err := sessionRepo.GetByToken(c.Request.Context(), parts[1])
        if err != nil || sess == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }

        c.Set("user_id", sess.UserID)
        c.Next()
    }
}
