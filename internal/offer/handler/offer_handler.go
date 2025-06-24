package handler

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	"github.com/gin-gonic/gin"

	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/service"
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

func (h *offerHandler) CreateOffer(c *gin.Context) {
	var req request.CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := session.Manager.GetInt(c.Request.Context(), "user_id")
	o, err := h.svc.Create(c.Request.Context(), userID, service.CreateOfferInput{
		ApplicationID: req.ApplicationID,
		Position:      req.Position,
		Salary:        req.Salary,
		StartDate:     req.StartDate,
		Location:      req.Location,
		Benefits:      req.Benefits,
		LetterContent: req.LetterContent,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ListOffersResponse{
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
	})
}

func (h *offerHandler) ListOffers(c *gin.Context) {
	var q request.ListOffersRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offers, total, err := h.svc.List(c.Request.Context(), q.Status, q.Page, q.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]response.OfferItem, len(offers))
	for i, o := range offers {
		items[i] = response.OfferItem{
			OfferID:       o.OfferID,
			ApplicationID: o.ApplicationID,
			Position:      o.Position,
			Salary:        o.Salary,
			StartDate:     o.StartDate,
			Location:      o.Location,
			Status:        o.Status,
		}
	}
	c.JSON(http.StatusOK, response.ListOffersResponse{Total: total, Items: items})
}

func (h *offerHandler) GetOfferDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("offer_id"))
	o, approvers, err := h.svc.GetDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, response.OfferDetail{
		OfferID:   o.OfferID,
		Approvers: approvers,
		Status:    o.Status,
	})
}

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
	c.JSON(http.StatusOK, response.SendOfferEmailResponse{
		EmailHistoryID: sent.EmailID,
		SentAt:         sent.CreatedAt,
	})
}
