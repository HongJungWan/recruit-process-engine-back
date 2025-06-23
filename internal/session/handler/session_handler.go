package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session/service"
)

type SessionHandler interface {
    Login(c *gin.Context)
    Logout(c *gin.Context)
}

type sessionHandler struct {
    svc service.SessionService
}

func NewSessionHandler(svc service.SessionService) SessionHandler {
    return &sessionHandler{svc: svc}
}

func (h *sessionHandler) Login(c *gin.Context) {
    var dto request.LoginRequest
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, token, err := h.svc.Login(c.Request.Context(), dto.LoginID, dto.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    resp := response.LoginResponse{
        SessionToken: token,
        User: response.UserInfo{
            UserID:    user.UserID,
            Email:     user.Email,
            Name:      user.Name,
            CreatedAt: user.CreatedAt,
            CreatedBy: user.CreatedBy,
        },
    }
    c.JSON(http.StatusOK, resp)
}

func (h *sessionHandler) Logout(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    parts := strings.SplitN(auth, " ", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
        return
    }
    if err := h.svc.Logout(c.Request.Context(), parts[1]); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
