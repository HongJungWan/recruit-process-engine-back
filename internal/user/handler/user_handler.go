package handler

import (
	"net/http"

	service "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
    HealthCheck(c *gin.Context)
}

type userHandler struct {
    userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
    return &userHandler{userSvc: userSvc}
}

func (h *userHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
