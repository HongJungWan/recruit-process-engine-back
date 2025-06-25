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

// ListApprovals godoc
// @Summary      오퍼별 승인 이력 조회
// @Description  오퍼의 승인(결재) 이력 목록을 반환합니다.
// @Tags         Approvals
// @Param        offer_id  path      int     true  "오퍼 식별자"
// @Success      200       {array}   response.ApprovalHistoryItem
// @Failure      500
// @Router       /offers/{offer_id}/approvals [get]
func (h *approvalHandler) ListApprovals(c *gin.Context) {
    offerID, _ := strconv.Atoi(c.Param("offer_id"))

    list, err := h.svc.ListByOffer(c.Request.Context(), offerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := make([]response.ApprovalHistoryItem, len(list))
    for i, ap := range list {
        output[i] = response.ApprovalHistoryItem{
            ApprovalID:  ap.ApprovalID,
            ApproverID:  ap.ApproverID,
            Status:      ap.Status,
            Comment:     ap.Comment,
            RequestedAt: ap.RequestedAt,
            DecidedAt:   ap.DecidedAt,
        }
    }
    c.JSON(http.StatusOK, output)
}

// CreateApprovals godoc
// @Summary      승인 요청 생성
// @Description  오퍼 승인(결재) 요청을 생성한다.
// @Tags         Approvals
// @Param        offer_id     path      int                              true  "오퍼 식별자"
// @Param        body         body      request.CreateApprovalsRequest   true  "승인자 ID 리스트"
// @Success      200          {array}   response.CreateApprovalsResponse
// @Failure      400
// @Failure      500
// @Router       /offers/{offer_id}/approvals [post]
func (h *approvalHandler) CreateApprovals(c *gin.Context) {
    offerID, _ := strconv.Atoi(c.Param("offer_id"))

    var input request.CreateApprovalsRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    created, err := h.svc.Request(c.Request.Context(), offerID, input.ApproverIDs, strconv.Itoa(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := make(response.CreateApprovalsResponse, len(created))
    for i, ap := range created {
        output[i] = response.ApprovalHistoryItem{
            ApprovalID:  ap.ApprovalID,
            ApproverID:  ap.ApproverID,
            Status:      ap.Status,
            RequestedAt: ap.RequestedAt,
        }
    }
    c.JSON(http.StatusOK, output)
}

// ProcessApproval godoc
// @Summary      승인 처리
// @Description  오퍼 승인(결재) 요청을 처리한다.
// @Tags         Approvals
// @Param        approval_id  path      int                              true  "승인 요청 식별자"
// @Param        body         body      request.ProcessApprovalRequest   true  "승인 상태 및 코멘트"
// @Success      200          {object}  response.ProcessApprovalResponse
// @Failure      400
// @Router       /offers/{offer_id}/approvals/{approval_id} [put]
func (h *approvalHandler) ProcessApproval(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("approval_id"))

    var input request.ProcessApprovalRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")

    comment := ""
    if input.Comment != nil {
        comment = *input.Comment
    }

    decidedAt, err := h.svc.Process(
        c.Request.Context(),
        id,
        input.Status,
        comment,
        strconv.Itoa(userID),
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    output := response.ProcessApprovalResponse{
        ApprovalID: id,
        Status:     input.Status,
        DecidedAt:  decidedAt,
    }
    c.JSON(http.StatusOK, output)
}
