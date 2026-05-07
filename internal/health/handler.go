package health

import (
	"database/sql"
	"net/http"

	"clinic-api/internal/shared"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	if err := h.db.Ping(); err != nil {
		shared.WriteError(w, http.StatusServiceUnavailable, "base de datos no disponible")
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "API y base de datos funcionando",
	})
}
