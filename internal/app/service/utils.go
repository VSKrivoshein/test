package service

import (
	e "github.com/VSKrivoshein/test/internal/app/custom_err"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", e.New(err, e.ErrPassword, http.StatusInternalServerError)

	}
	return string(bytes), nil
}
