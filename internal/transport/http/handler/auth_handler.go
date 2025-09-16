package handler

import (
	"net/http"

	"github.com/ariyaagustian/gin-boilerplate/internal/service"
	"github.com/ariyaagustian/gin-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(s service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

// Register godoc
// @Summary     Register user baru
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       payload body     dto.RegisterReq  true "Register payload"
// @Success     201     {object} dto.RegisterResp
// @Failure     400     {object} apperr.AppError
// @Router      /auth/register [post]
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

// Login godoc
// @Summary      Login dan dapatkan token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload body     dto.LoginReq true "Login payload"
// @Success      200     {object} dto.TokenResp
// @Failure      401     {object} apperr.AppError
// @Router       /auth/login [post]
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

// Me godoc
// @Summary      Get current user profile
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.User
// @Failure      401 {object} apperr.AppError
// @Router       /api/v1/users/me [get]
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

// AdminSetPassword godoc
// @Summary      Set password user (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body     map[string]string true "email, new_password"
// @Success      200     {object} map[string]string
// @Failure      403     {object} apperr.AppError
// @Router       /api/v1/admin/users/set-password [post]
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
