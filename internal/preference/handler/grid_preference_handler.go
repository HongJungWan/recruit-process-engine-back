package handler

import (
	// 표준 라이브러리
	"encoding/json"
	"net/http"
	"strconv"

	// 서드파티(외부) 라이브러리
	"github.com/gin-gonic/gin"

	// 내부 패키지
	req "github.com/HongJungWan/recruit-process-engine-back/internal/preference/dto/request"
	res "github.com/HongJungWan/recruit-process-engine-back/internal/preference/dto/response"
	svc "github.com/HongJungWan/recruit-process-engine-back/internal/preference/service"
	"github.com/HongJungWan/recruit-process-engine-back/internal/session"
)

type GridPreferenceHandler interface {
    GetGridPreferences(c *gin.Context)
    CreateGridPreference(c *gin.Context)
    UpdateGridPreference(c *gin.Context)
    DeleteGridPreference(c *gin.Context)
}

type gridPreferenceHandler struct {
    svc svc.GridPreferenceService
}

func NewGridPreferenceHandler(s svc.GridPreferenceService) GridPreferenceHandler {
    return &gridPreferenceHandler{svc: s}
}

func (h *gridPreferenceHandler) GetGridPreferences(c *gin.Context) {
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    if userID == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "존재하지 않는 유저입니다."})
        return
    }

    list, err := h.svc.GetByUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    out := make([]res.GridPreference, len(list))
    for i, gp := range list {
        json.Unmarshal(gp.Config, &out[i].Config)
        out[i].PreferenceID = gp.PreferenceID
        out[i].UserID = gp.UserID
        out[i].GridName = gp.GridName
        out[i].CreatedAt = gp.CreatedAt
        out[i].CreatedBy = gp.CreatedBy
        out[i].UpdatedAt = gp.UpdatedAt
        out[i].UpdatedBy = gp.UpdatedBy
    }
    c.JSON(http.StatusOK, out)
}

func (h *gridPreferenceHandler) CreateGridPreference(c *gin.Context) {
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    if userID == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "존재하지 않는 유저입니다."})
        return
    }

    var input req.CreateGridPreference
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    createdBy := strconv.Itoa(userID) // 또는 사용자 이름
    gp, err := h.svc.Create(c.Request.Context(), userID, input, createdBy)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var out res.GridPreference
    json.Unmarshal(gp.Config, &out.Config)
    out.PreferenceID = gp.PreferenceID
    out.UserID = gp.UserID
    out.GridName = gp.GridName
    out.CreatedAt = gp.CreatedAt
    out.CreatedBy = gp.CreatedBy
    c.JSON(http.StatusOK, out)
}

func (h *gridPreferenceHandler) UpdateGridPreference(c *gin.Context) {
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    if userID == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "존재하지 않는 유저입니다."})
        return
    }

    pid, _ := strconv.Atoi(c.Param("preference_id"))

    var input req.UpdateGridPreference
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    updatedBy := strconv.Itoa(userID)
    
    gp, err := h.svc.Update(c.Request.Context(), userID, pid, input, updatedBy)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var out res.GridPreference
    json.Unmarshal(gp.Config, &out.Config)
    out.PreferenceID = gp.PreferenceID
    out.UserID = gp.UserID
    out.GridName = gp.GridName
    out.UpdatedAt = gp.UpdatedAt
    out.UpdatedBy = gp.UpdatedBy
    c.JSON(http.StatusOK, out)
}

func (h *gridPreferenceHandler) DeleteGridPreference(c *gin.Context) {
    userID := session.Manager.GetInt(c.Request.Context(), "user_id")
    if userID == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "존재하지 않는 유저입니다."})
        return
    }

    pid, _ := strconv.Atoi(c.Param("preference_id"))

    if err := h.svc.Delete(c.Request.Context(), userID, pid); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
