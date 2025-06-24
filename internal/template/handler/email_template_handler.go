package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	req "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	"github.com/gin-gonic/gin"
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
    out := make([]res.EmailTemplateItem, len(list))
    for i, t := range list {
        var cfg map[string]interface{}
        json.Unmarshal(t.Config, &cfg)
        out[i] = res.EmailTemplateItem{
            TemplateID: t.ID,
            Name:       t.Name,
            Config:     cfg,
            CreatedAt:  t.CreatedAt,
        }
    }
    c.JSON(http.StatusOK, out)
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
    c.JSON(http.StatusOK, res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    })
}

func (h *emailTemplateHandler) CreateTemplate(c *gin.Context) {
    var reqBody req.CreateEmailTemplate
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    t, err := h.svc.Create(c.Request.Context(), reqBody.Name, reqBody.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    var cfg map[string]interface{}
    json.Unmarshal(t.Config, &cfg)
    c.JSON(http.StatusOK, res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    })
}

func (h *emailTemplateHandler) UpdateTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))
    var reqBody req.UpdateEmailTemplate
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    t, err := h.svc.Update(c.Request.Context(), id, reqBody.Name, reqBody.Config)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    var cfg map[string]interface{}
    json.Unmarshal(t.Config, &cfg)
    c.JSON(http.StatusOK, res.EmailTemplateDetail{
        TemplateID: t.ID,
        Name:       t.Name,
        Config:     cfg,
        CreatedAt:  t.CreatedAt,
    })
}

func (h *emailTemplateHandler) DeleteTemplate(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("template_id"))
    if err := h.svc.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
