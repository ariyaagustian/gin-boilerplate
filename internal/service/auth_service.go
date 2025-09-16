package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ariyaagustian/gin-boilerplate/internal/domain"
	"github.com/ariyaagustian/gin-boilerplate/internal/repository"
	"github.com/ariyaagustian/gin-boilerplate/pkg/apperr"
	"github.com/ariyaagustian/gin-boilerplate/pkg/auth"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, string, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	AdminSetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

type authSvc struct {
	repo      repository.UserRepository
	v         *validator.Validate
	jwtSecret string
	accessTTL time.Duration
}

func NewAuthSvc(r repository.UserRepository, v *validator.Validate, jwtSecret string, accessTTL time.Duration) AuthService {
	if v == nil {
		v = validator.New()
	}
	return &authSvc{repo: r, v: v, jwtSecret: jwtSecret, accessTTL: accessTTL}
}

type regDTO struct {
	Name     string `validate:"required,min=2"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

func (s *authSvc) Register(ctx context.Context, name, email, password string) (*domain.User, string, error) {
	in := regDTO{
		Name:     strings.TrimSpace(name),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: strings.TrimSpace(password),
	}
	if err := s.v.Struct(in); err != nil {
		return nil, "", apperr.Validation(err.Error(), err)
	}

	// cek email existing
	if _, err := s.repo.FindByEmail(ctx, in.Email); err == nil {
		return nil, "", apperr.Conflict("email sudah terdaftar", nil)
	}

	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", apperr.Internal("gagal hash password", err)
	}
	hs := string(hash)

	u := &domain.User{Name: in.Name, Email: in.Email, PasswordHash: &hs}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, "", apperr.Internal("gagal menyimpan user", err)
	}

	tok, _, err := auth.NewAccessToken(s.jwtSecret, u.ID, u.Email, s.accessTTL)
	if err != nil {
		return nil, "", apperr.Internal("gagal membuat token", err)
	}
	return u, tok, nil
}

// internal/service/auth_service.go (potongan)
func (s *authSvc) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if err := s.v.Var(email, "required,email"); err != nil {
		return nil, "", apperr.Validation("email tidak valid", err)
	}
	if password == "" {
		return nil, "", apperr.Validation("password wajib diisi", nil)
	}

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", apperr.Unauthorized("email atau password salah", err)
	}

	// ⬇️ Tambahan: jika user belum punya password
	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return nil, "", apperr.Unauthorized("akun belum memiliki password, minta admin untuk set password terlebih dahulu", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)); err != nil {
		return nil, "", apperr.Unauthorized("email atau password salah", err)
	}

	tok, _, err := auth.NewAccessToken(s.jwtSecret, u.ID, u.Email, s.accessTTL)
	if err != nil {
		return nil, "", apperr.Internal("gagal membuat token", err)
	}
	return u, tok, nil
}

func (s *authSvc) AdminSetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	newPassword = strings.TrimSpace(newPassword)
	if err := s.v.Var(newPassword, "required,min=6"); err != nil {
		return apperr.Validation("password minimal 6 karakter", err)
	}

	// pastikan user ada
	if _, err := s.repo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("user tidak ditemukan", err)
		}
		return apperr.Internal("gagal mengambil user", err)
	}

	// hash baru
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Internal("gagal hash password", err)
	}
	if err := s.repo.UpdatePasswordHash(ctx, userID, string(hash)); err != nil {
		return apperr.Internal("gagal menyimpan password", err)
	}
	return nil
}
