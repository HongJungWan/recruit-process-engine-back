package handler

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	req "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	"github.com/gin-gonic/gin"
)

type EmailHistoryHandler interface {
    SendEmail(c *gin.Context)
    ListHistory(c *gin.Context)
    GetHistory(c *gin.Context)
}

type emailHistoryHandler struct {
    svc service.EmailHistoryService
}

func NewEmailHistoryHandler(s svc.EmailHistoryService) EmailHistoryHandler {
    return &emailHistoryHandler{svc: s}
}

func (h *emailHistoryHandler) SendEmail(c *gin.Context) {
    var reqBody req.CreateEmailHistory
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    hist := &model.EmailHistory{
        UserID:        &userID,
        ApplicationID: reqBody.ApplicantID,
        OfferID:       reqBody.OfferID,
        TemplateID:    reqBody.TemplateID,
        Title:         reqBody.Title,
        Body:          reqBody.Body,
        CreatedBy:     strconv.Itoa(userID),
    }
    sent, err := h.svc.Send(c.Request.Context(), hist)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res.CreateEmailHistoryResponse{
        EmailID: sent.EmailID,
        SentAt:  sent.CreatedAt,
    })
}

func (h *emailHistoryHandler) ListHistory(c *gin.Context) {
    var q req.ListEmailHistory
    if err := c.ShouldBindQuery(&q); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    list, err := h.svc.List(c.Request.Context(), q.ApplicantID, q.OfferID, q.Page, q.Size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out := make([]res.EmailHistoryItem, len(list))
    for i, hst := range list {
        out[i] = res.EmailHistoryItem{
            EmailID: hst.EmailID,
            Title:   hst.Title,
            SentAt:  hst.CreatedAt,
        }
    }
    c.JSON(http.StatusOK, out)
}

func (h *emailHistoryHandler) GetHistory(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("email_id"))
    hst, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, res.EmailHistoryDetail{
        EmailID: hst.EmailID,
        Title:   hst.Title,
        Body:    hst.Body,
        SentAt:  hst.CreatedAt,
    })
}
