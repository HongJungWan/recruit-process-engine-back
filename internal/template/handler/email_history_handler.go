package handler

import (
	// 표준 라이브러리
	"net/http"
	"strconv"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	req "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/template/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
	"github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
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
    var input req.CreateEmailHistory
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    hist := &model.EmailHistory{
        UserID:        &userID,
        ApplicationID: input.ApplicantID,
        OfferID:       input.OfferID,
        TemplateID:    input.TemplateID,
        Title:         input.Title,
        Body:          input.Body,
        CreatedBy:     strconv.Itoa(userID),
    }

    sent, err := h.svc.Send(c.Request.Context(), hist)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := res.CreateEmailHistoryResponse{
        EmailID: sent.EmailID,
        SentAt:  sent.CreatedAt,
    }
    c.JSON(http.StatusOK, output)
}

func (h *emailHistoryHandler) ListHistory(c *gin.Context) {
    var input req.ListEmailHistory
    if err := c.ShouldBindQuery(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    list, err := h.svc.List(c.Request.Context(), input.ApplicantID, input.OfferID, input.Page, input.Size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := make([]res.EmailHistoryItem, len(list))
    for i, hst := range list {
        output[i] = res.EmailHistoryItem{
            EmailID: hst.EmailID,
            Title:   hst.Title,
            SentAt:  hst.CreatedAt,
        }
    }
    c.JSON(http.StatusOK, output)
}

func (h *emailHistoryHandler) GetHistory(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("email_id"))

    hst, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    output := res.EmailHistoryDetail{
        EmailID: hst.EmailID,
        Title:   hst.Title,
        Body:    hst.Body,
        SentAt:  hst.CreatedAt,
    }
    c.JSON(http.StatusOK, output)
}
