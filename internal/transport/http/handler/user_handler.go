package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ariyaagustian/gin-crud-boilerplate/internal/service"
	"github.com/ariyaagustian/gin-crud-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-crud-boilerplate/pkg/response"
)

type UserHandler struct{ svc service.UserService }

func NewUserHandler(s service.UserService) *UserHandler { return &UserHandler{svc: s} }

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

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, out)
}

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

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.WriteError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"deleted": true})
}
