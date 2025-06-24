package handler

import (
	"net/http"
	"strconv"

	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/applicant/dto/response"
	"github.com/HongJungWan/recruit-process-engine-back/internal/applicant/service"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
	"github.com/gin-gonic/gin"
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
    var q request.ListApplicantsRequest
    if err := c.ShouldBindQuery(&q); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    items, total, err := h.svc.List(c.Request.Context(), q.Page, q.Size, q.Stage, q.Keyword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out := res.ListApplicantsResponse{
        Total: total,
        Page:  q.Page,
        Size:  q.Size,
        Items: make([]res.ApplicantItem, len(items)),
    }
    for i, a := range items {
        out.Items[i] = res.ApplicantItem{
            ApplicationID: a.ApplicationID,
            Name:          a.Name,
            Email:         a.Email,
            CurrentStage:  a.CurrentStage,
        }
    }
    c.JSON(http.StatusOK, out)
}

func (h *applicantHandler) GetApplicant(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("application_id"))
    a, err := h.svc.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, res.ApplicantDetail{
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
    })
}

func (h *applicantHandler) UpdateApplicantStage(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("application_id"))
    var req request.UpdateStageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    old, updatedAt, err := h.svc.UpdateStage(c.Request.Context(), id, req.Stage, strconv.Itoa(userID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res.UpdateStageResponse{
        ApplicationID: id,
        OldStage:      old,
        NewStage:      req.Stage,
        UpdatedAt:     updatedAt,
    })
}

func (h *applicantHandler) BulkUpdateApplicantStage(c *gin.Context) {
    var req request.BulkUpdateStageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    cnt, err := h.svc.BulkUpdateStage(c.Request.Context(), req.IDs, req.Stage, strconv.Itoa(userID))
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
    out := make([]res.StageHistoryItem, len(hs))
    for i, hst := range hs {
        out[i] = res.StageHistoryItem{
            HistoryID: hst.HistoryID,
            Stage:     hst.Stage,
            Status:    hst.Status,
            CreatedAt: hst.CreatedAt,
            CreatedBy: hst.CreatedBy,
            UpdatedAt: hst.UpdatedAt,
            UpdatedBy: hst.UpdatedBy,
        }
    }
    c.JSON(http.StatusOK, out)
}
