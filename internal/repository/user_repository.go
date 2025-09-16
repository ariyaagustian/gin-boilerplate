package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ariyaagustian/gin-boilerplate/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindPaged(ctx context.Context, q string, page, pageSize int, sortBy, sortDir string) ([]domain.User, int64, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) FindAll(ctx context.Context) ([]domain.User, error) {
	var out []domain.User
	// hindari load semua data di prod; tapi kalau tetap mau:
	err := r.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Find(&out).Error
	return out, err
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

// =========================
// Paged + Search + Sort
// =========================

const MaxPageSize = 100

var allowedSort = map[string]struct{}{
	"created_at": {},
	"name":       {},
	"email":      {},
}

// scopes
func scopeSearch(q string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if s := strings.TrimSpace(q); s != "" {
			like := "%" + s + "%"
			// Use database-agnostic case-insensitive search
			// GORM will handle the appropriate SQL syntax for different databases
			return db.Where("LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)", like, like)
		}
		return db
	}
}

func scopeSort(sortBy, sortDir string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Validate and sanitize sort column
		if sortBy == "" {
			sortBy = "created_at"
		}

		// Use strict whitelist validation to prevent SQL injection
		if _, ok := allowedSort[sortBy]; !ok {
			sortBy = "created_at"
		}

		// Validate sort direction
		desc := true
		if strings.EqualFold(sortDir, "asc") {
			desc = false
		}

		return db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortBy},
			Desc:   desc,
		})
	}
}

func scopePaginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
		offset := (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *userRepo) FindPaged(
	ctx context.Context,
	q string,
	page, pageSize int,
	sortBy, sortDir string,
) ([]domain.User, int64, error) {

	var (
		items []domain.User
		total int64
	)

	base := r.db.WithContext(ctx).Model(&domain.User{})

	// hitung total
	if err := base.Scopes(scopeSearch(q)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// optional: kalau offset di luar total, kembalikan kosong cepat
	if total == 0 {
		return []domain.User{}, 0, nil
	}

	// ambil data
	if err := base.
		// pilih kolom yang diperlukan saja agar hemat memori (opsional)
		// Select("id", "name", "email", "phone", "created_at").
		Scopes(scopeSearch(q), scopeSort(sortBy, sortDir), scopePaginate(page, pageSize)).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"password_hash": hash,
			"updated_at":    gorm.Expr("now()"),
		}).Error
}
