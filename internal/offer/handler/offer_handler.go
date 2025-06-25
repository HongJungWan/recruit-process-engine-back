package handler

import (
	// 표준 라이브러리
	"net/http"
	"strconv"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/service"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	templateModel "github.com/HongJungWan/recruit-process-engine-back/internal/template/model"
	templateService "github.com/HongJungWan/recruit-process-engine-back/internal/template/service"
)

type OfferHandler interface {
    CreateOffer(c *gin.Context)
    ListOffers(c *gin.Context)
    GetOfferDetail(c *gin.Context)
    SendEmail(c *gin.Context)
}

type offerHandler struct {
    svc     service.OfferService
    mailSvc templateService.EmailHistoryService
}

func NewOfferHandler(svc service.OfferService, mailSvc templateService.EmailHistoryService) OfferHandler {
    return &offerHandler{svc: svc, mailSvc: mailSvc}
}

// CreateOffer godoc
// @Summary      오퍼 생성
// @Description  새로운 오퍼를 생성한다.
// @Tags         Offers
// @Param        body  body    request.CreateOfferRequest  true  "오퍼 생성 요청"
// @Success      200   {object}  response.ListOffersResponse
// @Failure      400
// @Failure      500
// @Router       /offers [post]
func (h *offerHandler) CreateOffer(c *gin.Context) {
    var input request.CreateOfferRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    o, err := h.svc.Create(c.Request.Context(), userID, service.CreateOfferInput{
        ApplicationID: input.ApplicationID,
        Position:      input.Position,
        Salary:        input.Salary,
        StartDate:     input.StartDate,
        Location:      input.Location,
        Benefits:      input.Benefits,
        LetterContent: input.LetterContent,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	output := response.ListOffersResponse{
        Total: 1,
        Items: []response.OfferItem{{
            OfferID:       o.OfferID,
            ApplicationID: o.ApplicationID,
            Position:      o.Position,
            Salary:        o.Salary,
            StartDate:     o.StartDate,
            Location:      o.Location,
            Status:        o.Status,
        }},
    }
    c.JSON(http.StatusOK, output)
}

// ListOffers godoc
// @Summary      오퍼 목록 조회
// @Description  오퍼 목록과 총 개수를 조회한다.
// @Tags         Offers
// @Param        status  query   string  false  "오퍼 상태 필터"
// @Param        page    query   int     false  "페이지 번호"
// @Param        size    query   int     false  "페이지당 항목 수"
// @Success      200     {object}  response.ListOffersResponse
// @Failure      400
// @Failure      500
// @Router       /offers [get]
func (h *offerHandler) ListOffers(c *gin.Context) {
    var input request.ListOffersRequest
    if err := c.ShouldBindQuery(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    offers, total, err := h.svc.List(c.Request.Context(), input.Status, input.Page, input.Size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := make([]response.OfferItem, len(offers))
    for i, o := range offers {
        output[i] = response.OfferItem{
            OfferID:       o.OfferID,
            ApplicationID: o.ApplicationID,
            Position:      o.Position,
            Salary:        o.Salary,
            StartDate:     o.StartDate,
            Location:      o.Location,
            Status:        o.Status,
        }
    }
    c.JSON(http.StatusOK, response.ListOffersResponse{Total: total, Items: output})
}

// GetOfferDetail godoc
// @Summary      오퍼 상세 조회
// @Description  오퍼의 상세 정보와 승인자 목록을 조회한다.
// @Tags         Offers
// @Param        offer_id  path      int  true  "오퍼 식별자"
// @Success      200       {object}  response.OfferDetail
// @Failure      404
// @Router       /offers/{offer_id} [get]
func (h *offerHandler) GetOfferDetail(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("offer_id"))

    o, approvers, err := h.svc.GetDetail(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

	output := response.OfferDetail{
        OfferID:   o.OfferID,
        Approvers: approvers,
        Status:    o.Status,
    }
    c.JSON(http.StatusOK, output)
}

// SendEmail godoc
// @Summary      오퍼 레터 이메일 전송
// @Description  오퍼 레터 이메일을 전송하고, 발송 이력을 기록하여 반환한다.
// @Tags         Offers
// @Param        offer_id  path      int  true  "오퍼 식별자"
// @Success      200       {object}  response.SendOfferEmailResponse
// @Failure      404
// @Failure      500
// @Router       /offers/{offer_id}/send [post]
func (h *offerHandler) SendEmail(c *gin.Context) {
    offerID, _ := strconv.Atoi(c.Param("offer_id"))
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")

    o, _, err := h.svc.GetDetail(c.Request.Context(), offerID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    hist := &templateModel.EmailHistory{
        UserID:        &userID,
        ApplicationID: &o.ApplicationID,
        OfferID:       &offerID,
        TemplateID:    0,
        Title:         "오퍼 레터: " + o.Position,
        Body:          o.LetterContent,
        CreatedBy:     strconv.Itoa(userID),
    }
    sent, err := h.mailSvc.Send(c.Request.Context(), hist)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	output := response.SendOfferEmailResponse{
        EmailHistoryID: sent.EmailID,
        SentAt:         sent.CreatedAt,
    }
    c.JSON(http.StatusOK, output)
}
