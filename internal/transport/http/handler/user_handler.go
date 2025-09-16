package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ariyaagustian/gin-boilerplate/internal/service"
	"github.com/ariyaagustian/gin-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-boilerplate/pkg/response"
)

type UserHandler struct{ svc service.UserService }

func NewUserHandler(s service.UserService) *UserHandler { return &UserHandler{svc: s} }

// Create godoc
// @Summary      Create user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload body     dto.CreateUserReq true "User payload"
// @Success      201     {object} domain.User
// @Failure      400     {object} apperr.AppError
// @Failure      401     {object} apperr.AppError
// @Router       /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.WriteError(c, apperr.BadRequest("payload tidak valid", err))
		return
	}
	out, err := h.svc.Create(c.Request.Context(), in.Name, in.Email)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": out})
}

// List godoc
// @Summary      List users
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        page   query    int false "page"   example(1)
// @Param        limit  query    int false "limit"  example(20)
// @Success      200    {array}  dto.ListUsersResp
// @Failure      401    {object} apperr.AppError
// @Router       /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Ambil query params → siapkan default
	p := service.ListUsersParams{
		Q:       c.Query("q"),
		SortBy:  c.DefaultQuery("sort_by", "created_at"),
		SortDir: c.DefaultQuery("sort_dir", "desc"),
	}

	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.Page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.PageSize = n
		}
	}

	out, err := h.svc.List(c.Request.Context(), p)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

// Get godoc
// @Summary      Get user by ID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "User ID (UUID)" format(uuid)
// @Success      200 {object} domain.User
// @Failure      404 {object} apperr.AppError
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

// Update godoc
// @Summary      Update user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path string       true "User ID (UUID)" format(uuid)
// @Param        payload body dto.UpdateUserReq true "Update payload"
// @Success      200     {object} domain.User
// @Failure      400     {object} apperr.AppError
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		response.WriteError(c, apperr.BadRequest("payload tidak valid", err))
		return
	}
	out, err := h.svc.Update(c.Request.Context(), id, in.Name, in.Email)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

// Delete godoc
// @Summary      Delete user
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "User ID (UUID)" format(uuid)
// @Success      204 {string} string "no content"
// @Failure      404 {object} apperr.AppError
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"deleted": true})
}
