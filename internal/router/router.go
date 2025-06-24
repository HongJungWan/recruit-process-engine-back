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

	appHand "github.com/HongJungWan/recruit-process-engine-back/internal/applicant/handler"
	appRepo "github.com/HongJungWan/recruit-process-engine-back/internal/applicant/repository"
	appSvc "github.com/HongJungWan/recruit-process-engine-back/internal/applicant/service"

	tmplHand "github.com/HongJungWan/recruit-process-engine-back/internal/template/handler"
	tmplRepo "github.com/HongJungWan/recruit-process-engine-back/internal/template/repository"
	tmplSvc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
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

    // 지원자 계층
    ar := appRepo.NewApplicantRepository(db)
    as := appSvc.NewApplicantService(ar)
    ah := appHand.NewApplicantHandler(as)

    // 이메일 템플릿 계층
    etRepo := tmplRepo.NewEmailTemplateRepository(db)
    etSvc  := tmplSvc.NewEmailTemplateService(etRepo)
    etH    := tmplHand.NewEmailTemplateHandler(etSvc)

    // 이메일 발송 이력 계층
    ehRepo := tmplRepo.NewEmailHistoryRepository(db)
    ehSvc  := tmplSvc.NewEmailHistoryService(ehRepo)
    ehH    := tmplHand.NewEmailHistoryHandler(ehSvc)

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

    api.GET("/applicants", ah.ListApplicants)
    api.GET("/applicants/:application_id", ah.GetApplicant)
    api.PATCH("/applicants/:application_id/stage", ah.UpdateApplicantStage)
    api.POST( "/applicants/stages/bulk-update", ah.BulkUpdateApplicantStage)
    api.GET("/applicants/:application_id/history", ah.GetApplicantHistory)

    api.GET("/email-templates",                   etH.ListTemplates)
    api.GET("/email-templates/:template_id",      etH.GetTemplate)
    api.POST("/email-templates",                  etH.CreateTemplate)
    api.PUT("/email-templates/:template_id",      etH.UpdateTemplate)
    api.DELETE("/email-templates/:template_id",   etH.DeleteTemplate)

    api.POST("/email-history",                    ehH.SendEmail)
    api.GET("/email-history",                     ehH.ListHistory)
    api.GET("/email-history/:email_id",           ehH.GetHistory)

    return r
}
