package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"clinic-api/internal/shared"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo       *Repository
	jwtManager *JWTManager
}

func NewHandler(repo *Repository, jwtManager *JWTManager) *Handler {
	return &Handler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	var input LoginRequest

	if err := shared.ReadJSON(w, r, &input); err != nil {
		log.Printf("error leyendo json login: %v", err)
		shared.WriteError(w, http.StatusBadRequest, "json invalido")
		return
	}

	if err := input.Validar(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.FindByUsername(input.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteError(w, http.StatusUnauthorized, "credenciales invalidas")
			return
		}
		log.Printf("error consultando usuario login: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al iniciar sesion")
		return
	}

	if !user.Activo {
		shared.WriteError(w, http.StatusUnauthorized, "usuario inactivo")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		shared.WriteError(w, http.StatusUnauthorized, "credenciales invalidas")
		return
	}

	token, err := h.jwtManager.Generate(user)
	if err != nil {
		log.Printf("error generando jwt: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al generar token")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
