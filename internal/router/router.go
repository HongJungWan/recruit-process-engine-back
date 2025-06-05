package router

import (
    "github.com/gin-gonic/gin"
    "github.com/HongJungWan/recruit-process-engine-back/internal/user/handler"
)

func InitRouter(userHandler handler.UserHandler) *gin.Engine {
    r := gin.Default()

    r.GET("/health", userHandler.HealthCheck)

    api := r.Group("/api")
    {
        users := api.Group("/users")
        {
            users.POST("", userHandler.Register)
            users.GET("/:id", userHandler.GetProfile)
        }

        auth := api.Group("/auth")
        {
            auth.POST("", userHandler.Login)
        }
    }

    return r
}
