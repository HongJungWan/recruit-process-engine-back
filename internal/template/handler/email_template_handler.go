package handler

import (
	// 표준 라이브러리
	"encoding/json"
	"net/http"
	"strconv"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	req "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
)

type EmailTemplateHandler interface {
    ListTemplates(c *gin.Context)
    GetTemplate(c *gin.Context)
    CreateTemplate(c *gin.Context)
    UpdateTemplate(c *gin.Context)
    DeleteTemplate(c *gin.Context)
}

type emailTemplateHandler struct {
    svc service.EmailTemplateService
}

func NewEmailTemplateHandler(s svc.EmailTemplateService) EmailTemplateHandler {
    return &emailTemplateHandler{svc: s}
}

// ListTemplates godoc
// @Summary     이메일 템플릿 목록 조회
// @Description 저장된 모든 이메일 템플릿의 정보를 반환한다.
// @Tags        EmailTemplates
// @Produce     json
// @Success     200  {array}   res.EmailTemplateItem
// @Failure     500
// @Router      /email-templates [get]
func (h *emailTemplateHandler) ListTemplates(c *gin.Context) {
    list, err := h.svc.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := make([]res.EmailTemplateItem, len(list))
    for i, t := range list {
        var cfg map[string]interface{}
        json.Unmarshal(t.Config, &cfg)
        output[i] = res.EmailTemplateItem{
            TemplateID: t.ID,
            Name:       t.Name,
            Config:     cfg,
            CreatedAt:  t.CreatedAt,
        }
    }
    c.JSON(http.StatusOK, output)
}

// GetTemplate godoc
// @Summary     이메일 템플릿 상세 조회
// @Description 선택된 이메일 템플릿의 상세 정보를 반환한다.
// @Tags        EmailTemplates
// @Produce     json
// @Param       template_id  path      int  true  "템플릿 식별자"
// @Success     200          {object}  res.EmailTemplateDetail
// @Failure     404
// @Router      /email-templates/{template_id} [get]
func (h *emailTemplateHandler) GetTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))

    t, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    var cfg map[string]interface{}
    json.Unmarshal(t.Config, &cfg)

    output := res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    }

    c.JSON(http.StatusOK, output)
}

// CreateTemplate godoc
// @Summary     이메일 템플릿 생성
// @Description 새로운 이메일 템플릿을 생성한다.
// @Tags        EmailTemplates
// @Accept      json
// @Produce     json
// @Param       body  body      req.CreateEmailTemplate  true  "생성할 템플릿 정보"
// @Success     200   {object}  res.EmailTemplateDetail
// @Failure     400
// @Failure     500
// @Router      /email-templates [post]
func (h *emailTemplateHandler) CreateTemplate(c *gin.Context) {
    var input req.CreateEmailTemplate
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    t, err := h.svc.Create(c.Request.Context(), input.Name, input.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var cfg map[string]interface{}
    json.Unmarshal(t.Config, &cfg)

    output := res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    }

    c.JSON(http.StatusOK, output)
}

// UpdateTemplate godoc
// @Summary     이메일 템플릿 수정
// @Description 이메일 템플릿 설정을 업데이트한다.
// @Tags        EmailTemplates
// @Accept      json
// @Produce     json
// @Param       template_id  path      int                          true  "템플릿 식별자"
// @Param       body         body      req.UpdateEmailTemplate      true  "업데이트할 템플릿 정보"
// @Success     200          {object}  res.EmailTemplateDetail
// @Failure     400
// @Failure     500
// @Router      /email-templates/{template_id} [put]
func (h *emailTemplateHandler) UpdateTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))

    var input req.UpdateEmailTemplate
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    t, err := h.svc.Update(c.Request.Context(), id, input.Name, input.Config)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var cfg map[string]interface{}
    json.Unmarshal(t.Config, &cfg)

    output := res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    }

    c.JSON(http.StatusOK, output)
}

// DeleteTemplate godoc
// @Summary     이메일 템플릿 삭제
// @Description 선택된 이메일 템플릿을 삭제한다.
// @Tags        EmailTemplates
// @Param       template_id  path  int  true  "템플릿 식별자"
// @Success     204
// @Failure     400
// @Failure     500
// @Router      /email-templates/{template_id} [delete]
func (h *emailTemplateHandler) DeleteTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))
    
    if err := h.svc.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
