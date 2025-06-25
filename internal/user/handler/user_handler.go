package handler

import (
	// 표준 라이브러리
	"net/http"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	"github.com/HongJungWan/recruit-process-engine-back/internal/user/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/user/service"
)

type UserHandler interface {
    HealthCheck(c *gin.Context)
    Login(c *gin.Context)
    Logout(c *gin.Context)
}

type userHandler struct {
    userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
    return &userHandler{userSvc: userSvc}
}

// HealthCheck godoc
// @Summary     서비스 헬스체크
// @Description 서비스 상태를 반환한다.
// @Tags        Health
// @Success     200
// @Router      /health-check [get]
func (h *userHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// Login godoc
// @Summary     사용자 로그인
// @Description 로그인 아이디와 비밀번호로 인증하고 세션을 생성한다.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credentials  body      request.Credentials  true  "로그인 정보"
// @Success     200          {object}  map[string]interface{}  "message와 user_id 반환"
// @Failure     400
// @Failure     401
// @Router      /auth/login [post]
func (h *userHandler) Login(c *gin.Context) {
    var input request.Credentials
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    userID, err := h.userSvc.Authenticate(c.Request.Context(), input.LoginId, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    session.Adapter.Put(c, "user_id", userID)
    c.JSON(http.StatusOK, gin.H{"message": "logged in", "user_id": userID})
}

// Logout godoc
// @Summary     사용자 로그아웃
// @Description 현재 세션을 파괴하여 로그아웃한다.
// @Tags        Auth
// @Success     204
// @Router      /auth/logout [post]
func (h *userHandler) Logout(c *gin.Context) {
    session.Adapter.Destroy(c)
    c.Status(http.StatusNoContent)
}
