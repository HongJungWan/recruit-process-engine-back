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
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
)

type EmailHistoryHandler interface {
    SendEmail(c *gin.Context)
    ListHistory(c *gin.Context)
    GetHistory(c *gin.Context)
}

type emailHistoryHandler struct {
    svc svc.EmailHistoryService
}

func NewEmailHistoryHandler(s svc.EmailHistoryService) EmailHistoryHandler {
    return &emailHistoryHandler{svc: s}
}

// SendEmail godoc
// @Summary     이메일 발송 이력 생성
// @Description 이메일 발송 이력 레코드를 생성한다.
// @Tags        EmailHistory
// @Accept      json
// @Produce     json
// @Param       body  body      req.CreateEmailHistory  true  "이메일 발송 요청"
// @Success     200   {object}  res.CreateEmailHistoryResponse
// @Failure     400
// @Failure     500
// @Router      /email-history [post]
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

// ListHistory godoc
// @Summary     이메일 발송 이력 목록 조회
// @Description 이메일 발송 이력 목록과 페이징 정보를 반환한다.
// @Tags        EmailHistory
// @Produce     json
// @Param       applicant_id  query     int     false  "지원자 식별자"
// @Param       offer_id      query     int     false  "오퍼 식별자"
// @Param       page          query     int     false  "페이지 번호"
// @Param       size          query     int     false  "페이지당 목록 수"
// @Success     200           {array}   res.EmailHistoryItem
// @Failure     400
// @Failure     500
// @Router      /email-history [get]
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

// GetHistory godoc
// @Summary     이메일 발송 이력 상세 조회
// @Description 특정 이메일 발송 이력의 상세 정보를 반환한다.
// @Tags        EmailHistory
// @Produce     json
// @Param       email_id  path      int  true  "이메일 이력 식별자"
// @Success     200       {object}  res.EmailHistoryDetail
// @Failure     404
// @Router      /email-history/{email_id} [get]
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
