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

func (h *emailTemplateHandler) DeleteTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))
    
    if err := h.svc.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
