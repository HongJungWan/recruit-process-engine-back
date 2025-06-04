package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    service "github.com/HongJungWan/recruit-process-engine-back/internal/user/service"
)

type UserHandler interface {
    Register(c *gin.Context)
    Login(c *gin.Context)
    GetProfile(c *gin.Context)
}

type userHandler struct {
    userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
    return &userHandler{userSvc: userSvc}
}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
}

func (h *userHandler) Register(c *gin.Context) {
    var req RegisterRequest
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

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *userHandler) Login(c *gin.Context) {
    var req LoginRequest

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
