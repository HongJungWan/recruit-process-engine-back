package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	userHandler "github.com/HongJungWan/recruit-process-engine-back/internal/user/handler"
	userRepo "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
	userSvc "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"

	sessionHandler "github.com/HongJungWan/recruit-process-engine-back/internal/session/handler"
	sessionRepo "github.com/HongJungWan/recruit-process-engine-back/internal/session/repository"
	sessionSvc "github.com/HongJungWan/recruit-process-engine-back/internal/session/service"

	"github.com/HongJungWan/recruit-process-engine-back/internal/middleware"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    // user layer
    ur := userRepo.NewUserRepository(db)
    us := userSvc.NewUserService(ur)
    uh := userHandler.NewUserHandler(us)

    // session layer
    sr := sessionRepo.NewSessionRepository(db)
    as := sessionSvc.NewSessionService(ur, sr, time.Hour*24) // TTL(만료 시간) 24h
    ah := sessionHandler.NewSessionHandler(as)

    // Public
    api := r.Group("/api/v1")
    _ = api.GET("/health-check", uh.HealthCheck)

    a := api.Group("/auth")
    {
        a.POST("/login", ah.Login)
        a.DELETE("/logout", middleware.Auth(sr), ah.Logout)
    }

    // TODO: Protected example

    return r
}
