package router

import (
    "github.com/gin-gonic/gin"
    "github.com/HongJungWan/recruit-process-engine-back/internal/handler"
)

func InitRouter(userHandler handler.UserHandler) *gin.Engine {
    r := gin.Default()

    api := r.Group("/api")
    {
        users := api.Group("/users")
        {
            users.POST("", userHandler.Register)
            users.POST("/login", userHandler.Login)
            users.GET("/:id", userHandler.GetProfile)
        }
    }

    return r
}
