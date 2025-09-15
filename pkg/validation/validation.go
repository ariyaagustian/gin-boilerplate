package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError menerima error dari validator.Struct()
// dan mengembalikan pesan yang lebih human-friendly.
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		// Ambil error pertama (bisa juga di-loop semua)
		fe := ve[0]
		switch fe.Tag() {
		case "required":
			return fmt.Sprintf("%s wajib diisi", strings.ToLower(fe.Field()))
		case "email":
			return "format email tidak valid"
		case "min":
			return fmt.Sprintf("%s minimal %s karakter", strings.ToLower(fe.Field()), fe.Param())
		default:
			return fmt.Sprintf("%s tidak valid", strings.ToLower(fe.Field()))
		}
	}
	return err.Error()
}
