package handler

import (
	"net/http"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	creds "github.com/HongJungWan/recruit-process-engine-back/internal/user/dto/request"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
    HealthCheck(c *gin.Context)
    Login(c *gin.Context)
    Logout(c *gin.Context)
}

type userHandler struct {
    userSvc svc.UserService
}

func NewUserHandler(userSvc svc.UserService) UserHandler {
    return &userHandler{userSvc: userSvc}
}

func (h *userHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *userHandler) Login(c *gin.Context) {
    var input creds.Credentials
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    userID, err := h.userSvc.Authenticate(c.Request.Context(), input.Username, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    session.Adapter.Put(c, "user_id", userID)
    c.JSON(http.StatusOK, gin.H{"message": "logged in", "user_id":  userID,})
}

func (h *userHandler) Logout(c *gin.Context) {
    session.Adapter.Destroy(c)
    c.Status(http.StatusNoContent)
}