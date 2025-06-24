package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	userHandler "github.com/HongJungWan/recruit-process-engine-back/internal/user/handler"
	userRepo "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
	userSvc "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    // 미들웨어 등록
    r.Use(session.Adapter.LoadAndSave)

    // user 계층
    ur := userRepo.NewUserRepository(db)
    us := userSvc.NewUserService(ur)
    uh := userHandler.NewUserHandler(us)

    api := r.Group("/api/v1")
    api.GET("/health-check", uh.HealthCheck)
    api.POST("/auth/login", uh.Login)
    api.POST("/auth/logout", uh.Logout)

    return r
}
