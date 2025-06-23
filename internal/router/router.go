package router

import (
	"github.com/HongJungWan/recruit-process-engine-back/internal/user/handler"
	"github.com/gin-gonic/gin"
)

func InitRouter(userHandler handler.UserHandler) *gin.Engine {
    r := gin.Default()

    r.GET("/health-check", userHandler.HealthCheck)

    return r
}
