package apperr

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

// AppError = error kustom dengan kode, pesan, dan HTTP status
type AppError struct {
	Code       string // ex: "duplicate", "validation", "bad_request"
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
	}
	return e.Message
}
func (e *AppError) Unwrap() error { return e.Err }

// ---------- Helper ctor ----------
func New(code string, httpStatus int, msg string, err error) *AppError {
	return &AppError{Code: code, HTTPStatus: httpStatus, Message: msg, Err: err}
}

func Validation(msg string, err error) *AppError   { return New("validation", 400, msg, err) }
func BadRequest(msg string, err error) *AppError   { return New("bad_request", 400, msg, err) }
func Conflict(msg string, err error) *AppError     { return New("duplicate", 409, msg, err) }
func NotFound(msg string, err error) *AppError     { return New("not_found", 404, msg, err) }
func Internal(msg string, err error) *AppError     { return New("internal", 500, msg, err) }
func Unauthorized(msg string, err error) *AppError { return New("unauthorized", 401, msg, err) }
func Forbidden(msg string, err error) *AppError    { return New("forbidden", 403, msg, err) }

// ---------- Parser khusus Postgres ----------
func FromPg(err error) *AppError {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case "23505": // unique_violation
		return Conflict("data sudah ada (duplikat)", err)
	case "23503": // foreign_key_violation
		return BadRequest("data masih direferensikan (foreign key)", err)
	case "23502": // not_null_violation
		return BadRequest("field wajib diisi (not null)", err)
		// di switch pgErr.Code tambahkan:
	case "22P02": // invalid_text_representation (e.g., UUID invalid)
		return BadRequest("format input tidak valid", err)
	case "22001": // string_data_right_truncation
		return BadRequest("panjang data melebihi batas", err)
	case "23514": // check_violation
		return BadRequest("data melanggar aturan (check constraint)", err)
	default:
		return BadRequest(pgErr.Message, err)
	}
}
