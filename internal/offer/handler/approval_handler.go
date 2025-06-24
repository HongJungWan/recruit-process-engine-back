package handler

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/request"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/offer/service"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	"github.com/gin-gonic/gin"
)

type ApprovalHandler interface {
    ListApprovals(c *gin.Context)
    CreateApprovals(c *gin.Context)
    ProcessApproval(c *gin.Context)
}

type approvalHandler struct {
    svc service.ApprovalService
}

func NewApprovalHandler(svc service.ApprovalService) ApprovalHandler {
    return &approvalHandler{svc: svc}
}

func (h *approvalHandler) ListApprovals(c *gin.Context) {
    offerID, _ := strconv.Atoi(c.Param("offer_id"))
    list, err := h.svc.ListByOffer(c.Request.Context(), offerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out := make([]response.ApprovalHistoryItem, len(list))
    for i, ap := range list {
        out[i] = response.ApprovalHistoryItem{
            ApprovalID:  ap.ApprovalID,
            ApproverID:  ap.ApproverID,
            Status:      ap.Status,
            Comment:     ap.Comment,
            RequestedAt: ap.RequestedAt,
            DecidedAt:   ap.DecidedAt,
        }
    }
    c.JSON(http.StatusOK, out)
}

func (h *approvalHandler) CreateApprovals(c *gin.Context) {
    offerID, _ := strconv.Atoi(c.Param("offer_id"))
    var reqBody request.CreateApprovalsRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    created, err := h.svc.Request(c.Request.Context(), offerID, reqBody.ApproverIDs, strconv.Itoa(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out := make(response.CreateApprovalsResponse, len(created))
    for i, ap := range created {
        out[i] = response.ApprovalHistoryItem{
            ApprovalID:  ap.ApprovalID,
            ApproverID:  ap.ApproverID,
            Status:      ap.Status,
            RequestedAt: ap.RequestedAt,
        }
    }
    c.JSON(http.StatusOK, out)
}

func (h *approvalHandler) ProcessApproval(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("approval_id"))
    var reqBody request.ProcessApprovalRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")

    comment := ""
    if reqBody.Comment != nil {
        comment = *reqBody.Comment
    }

    decidedAt, err := h.svc.Process(
        c.Request.Context(),
        id,
        reqBody.Status,
        comment,
        strconv.Itoa(userID),
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, response.ProcessApprovalResponse{
        ApprovalID: id,
        Status:     reqBody.Status,
        DecidedAt:  decidedAt,
    })
}
