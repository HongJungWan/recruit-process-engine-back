package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    reqdto "github.com/HongJungWan/recruit-process-engine-back/internal/user/dto/request"
    service "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"
)

type UserHandler interface {
    Register(c *gin.Context)
    Login(c *gin.Context)
    GetProfile(c *gin.Context)
    HealthCheck(c *gin.Context)
}

type userHandler struct {
    userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
    return &userHandler{userSvc: userSvc}
}

func (h *userHandler) Register(c *gin.Context) {
    var req reqdto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newID, err := h.userSvc.Register(c.Request.Context(), req.Email, req.Password, req.Name)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"id": newID})
}

func (h *userHandler) Login(c *gin.Context) {
    var req reqdto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userSvc.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "name": user.Name})
}

func (h *userHandler) GetProfile(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := h.userSvc.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *userHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
