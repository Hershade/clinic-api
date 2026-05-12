package auth

import (
	"errors"
	"strings"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Activo       bool   `json:"activo"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validar() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" {
		return errors.New("el username es obligatorio")
	}

	if r.Password == "" {
		return errors.New("el password es obligatorio")
	}

	return nil
}
