package handler

import (
	// 표준 라이브러리
	"net/http"
	"strconv"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/applicant/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/service"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
)

type ApplicantHandler interface {
    ListApplicants(c *gin.Context)
    GetApplicant(c *gin.Context)
    UpdateApplicantStage(c *gin.Context)
    BulkUpdateApplicantStage(c *gin.Context)
    GetApplicantHistory(c *gin.Context)
}

type applicantHandler struct {
    svc service.ApplicantService
}

func NewApplicantHandler(svc service.ApplicantService) ApplicantHandler {
    return &applicantHandler{svc: svc}
}

func (h *applicantHandler) ListApplicants(c *gin.Context) {
    var input request.ListApplicantsRequest
    if err := c.ShouldBindQuery(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    items, total, err := h.svc.List(c.Request.Context(), input.Page, input.Size, input.Stage, input.Keyword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    output := res.ListApplicantsResponse{
        Total: total,
        Page:  input.Page,
        Size:  input.Size,
        Items: make([]res.ApplicantItem, len(items)),
    }

    for i, a := range items {
        output.Items[i] = res.ApplicantItem{
            ApplicationID: a.ApplicationID,
            Name:          a.Name,
            Email:         a.Email,
            CurrentStage:  a.CurrentStage,
        }
    }
    c.JSON(http.StatusOK, output)
}

func (h *applicantHandler) GetApplicant(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("application_id"))

    a, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    output := res.ApplicantDetail{
        ApplicationID: a.ApplicationID,
        Name:          a.Name,
        Email:         a.Email,
        Phone:         a.Phone,
        Education:     a.Education,
        Experience:    a.Experience,
        TechStack:     a.TechStack,
        CurrentStage:  a.CurrentStage,
        CreatedAt:     a.CreatedAt,
        CreatedBy:     a.CreatedBy,
        UpdatedAt:     a.UpdatedAt,
        UpdatedBy:     a.UpdatedBy,
    }
    c.JSON(http.StatusOK, output)
}

func (h *applicantHandler) UpdateApplicantStage(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("application_id"))

    var input request.UpdateStageRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    
    old, updatedAt, err := h.svc.UpdateStage(c.Request.Context(), id, input.Stage, strconv.Itoa(userID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    output := res.UpdateStageResponse{
        ApplicationID: id,
        OldStage:      old,
        NewStage:      input.Stage,
        UpdatedAt:     updatedAt,
    }

    c.JSON(http.StatusOK, output)
}

func (h *applicantHandler) BulkUpdateApplicantStage(c *gin.Context) {
    var input request.BulkUpdateStageRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    
    cnt, err := h.svc.BulkUpdateStage(c.Request.Context(), input.IDs, input.Stage, strconv.Itoa(userID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res.BulkUpdateResponse{Updated: cnt})
}

func (h *applicantHandler) GetApplicantHistory(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("application_id"))

    hs, err := h.svc.GetHistory(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    output := make([]res.StageHistoryItem, len(hs))
    for i, hst := range hs {
        output[i] = res.StageHistoryItem{
            HistoryID: hst.HistoryID,
            Stage:     hst.Stage,
            Status:    hst.Status,
            CreatedAt: hst.CreatedAt,
            CreatedBy: hst.CreatedBy,
            UpdatedAt: hst.UpdatedAt,
            UpdatedBy: hst.UpdatedBy,
        }
    }
    c.JSON(http.StatusOK, output)
}
