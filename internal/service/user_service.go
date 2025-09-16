package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ariyaagustian/gin-boilerplate/internal/domain"
	"github.com/ariyaagustian/gin-boilerplate/internal/repository"
	"github.com/ariyaagustian/gin-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-boilerplate/pkg/validation"
)

const MaxPageSize = 100

type UserService interface {
	Create(ctx context.Context, name, email string) (*domain.User, error)
	List(ctx context.Context, p ListUsersParams) (PageResult[domain.User], error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id, name, email string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type userSvc struct {
	repo repository.UserRepository
	v    *validator.Validate
}

func NewUserSvc(r repository.UserRepository, v *validator.Validate) *userSvc {
	if v == nil {
		v = validator.New()
	}
	return &userSvc{repo: r, v: v}
}

type createUserDTO struct {
	Name  string `validate:"required,min=2"`
	Email string `validate:"required,email"`
}

type ListUsersParams struct {
	Q        string
	Page     int
	PageSize int
	SortBy   string // "created_at","name","email"
	SortDir  string // "asc" / "desc"
}

type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
}

func (s *userSvc) Create(ctx context.Context, name, email string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	dto := createUserDTO{Name: name, Email: email}
	if err := s.v.Struct(dto); err != nil {
		return nil, apperr.Validation(validation.FormatValidationError(err), err)
	}

	u := &domain.User{Name: dto.Name, Email: dto.Email}
	if err := s.repo.Create(ctx, u); err != nil {
		// 1) Heuristik pesan duplicate
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, apperr.Conflict("email sudah terdaftar, gunakan email lain", err)
		}
		// 2) Mapping kode PG
		if ae := apperr.FromPg(err); ae != nil {
			return nil, ae
		}
		// 3) Lainnya
		return nil, apperr.Internal("gagal menyimpan data", err)
	}
	return u, nil
}

func (s *userSvc) List(ctx context.Context, p ListUsersParams) (PageResult[domain.User], error) {
	// Normalisasi pagination
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	} else if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}

	// Default sorting
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortDir != "asc" && p.SortDir != "desc" {
		p.SortDir = "desc"
	}

	items, total, err := s.repo.FindPaged(ctx, p.Q, p.Page, p.PageSize, p.SortBy, p.SortDir)
	if err != nil {
		return PageResult[domain.User]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.PageSize)))
	return PageResult[domain.User]{
		Items:      items,
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
	}, nil
}

func (s *userSvc) Get(ctx context.Context, id string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.BadRequest("id tidak valid", err)
	}
	u, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("user tidak ditemukan", err)
		}
		return nil, apperr.Internal("gagal mengambil data", err)
	}
	return u, nil
}

func (s *userSvc) Update(ctx context.Context, id, name, email string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.BadRequest("id tidak valid", err)
	}

	u, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("user tidak ditemukan", err)
		}
		return nil, apperr.Internal("gagal mengambil data", err)
	}

	if name != "" {
		u.Name = strings.TrimSpace(name)
	}
	if email != "" {
		u.Email = strings.ToLower(strings.TrimSpace(email))
		// validasi email basic
		if err := s.v.Var(u.Email, "required,email"); err != nil {
			return nil, apperr.Validation(validation.FormatValidationError(err), err)
		}
	}

	if err := s.repo.Update(ctx, u); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, apperr.Conflict("email sudah terdaftar, gunakan email lain", err)
		}
		if ae := apperr.FromPg(err); ae != nil {
			return nil, ae
		}
		return nil, apperr.Internal("gagal menyimpan data", err)
	}
	return u, nil
}

func (s *userSvc) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperr.BadRequest("id tidak valid", err)
	}
	if err := s.repo.Delete(ctx, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("user tidak ditemukan", err)
		}
		return apperr.Internal("gagal menghapus data", err)
	}
	return nil
}
