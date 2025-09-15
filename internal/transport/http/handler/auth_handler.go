package handler

import (
	"net/http"

	"github.com/ariyaagustian/gin-crud-boilerplate/internal/service"
	"github.com/ariyaagustian/gin-crud-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-crud-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(s service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

func (h *AuthHandler) Register(c *gin.Context) {
	var in struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.WriteError(c, apperr.BadRequest("payload tidak valid", err))
		return
	}
	u, tok, err := h.svc.Register(c.Request.Context(), in.Name, in.Email, in.Password)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u, "token": tok, "token_type": "Bearer"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.WriteError(c, apperr.BadRequest("payload tidak valid", err))
		return
	}
	u, tok, err := h.svc.Login(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": u, "token": tok, "token_type": "Bearer"})
}

func (h *UserHandler) Me(c *gin.Context) {
	uid, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// reuse Get by ID
	out, err := h.svc.Get(c.Request.Context(), uid.(uuid.UUID).String())
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

func (h *AuthHandler) AdminSetPassword(c *gin.Context) {
	var in struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.WriteError(c, apperr.BadRequest("payload tidak valid", err))
		return
	}
	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		response.WriteError(c, apperr.BadRequest("user_id tidak valid", err))
		return
	}
	if err := h.svc.AdminSetPassword(c.Request.Context(), uid, in.Password); err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
