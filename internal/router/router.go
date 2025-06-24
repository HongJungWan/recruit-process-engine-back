package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session"

	userHandler "github.com/HongJungWan/recruit-process-engine-back/internal/user/handler"
	userRepo "github.com/HongJungWan/recruit-process-engine-back/internal/user/repository"
	userSvc "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"

	gpHand "github.com/HongJungWan/recruit-process-engine-back/internal/preference/handler"
	gpRepo "github.com/HongJungWan/recruit-process-engine-back/internal/preference/repository"
	gpSvc "github.com/HongJungWan/recruit-process-engine-back/internal/preference/service"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    // 미들웨어 등록
    r.Use(session.Adapter.LoadAndSave)

    // 유저 계층
    ur := userRepo.NewUserRepository(db)
    us := userSvc.NewUserService(ur)
    uh := userHandler.NewUserHandler(us)

    // 그리드 설정 계층
    gpRepo := gpRepo.NewGridPreferenceRepository(db)
    gpSvc  := gpSvc.NewGridPreferenceService(gpRepo)
    gpHand := gpHand.NewGridPreferenceHandler(gpSvc)

    api := r.Group("/api/v1")
    api.GET("/health-check", uh.HealthCheck)
    api.POST("/auth/login", uh.Login)
    api.POST("/auth/logout", uh.Logout)

    // TODO: 인증 받은 유저만 접근 가능하게 수정하기
    users := api.Group("/users")
    users.GET("/grid-preferences", gpHand.GetGridPreferences)
    users.POST("/grid-preferences", gpHand.CreateGridPreference)
    users.PUT("/grid-preferences/:preference_id", gpHand.UpdateGridPreference)
    users.DELETE("/grid-preferences/:preference_id", gpHand.DeleteGridPreference)

    return r
}
