package doctor

import (
	"clinic-api/internal/shared"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) DoctorsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDoctors(w, r)
	case http.MethodPost:
		h.createDoctor(w, r)
	default:
		shared.WriteError(w, http.StatusMethodNotAllowed, "Medoto no permitido")
	}
}

func (h *Handler) DoctorByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "Medoto no permitido")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/doctors/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "Id invalido")
	}

	doctor, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteError(w, http.StatusNotFound, "Doctor no encontrado")
			return
		}
		shared.WriteError(w, http.StatusInternalServerError, "Error al consultar doctor")
		return
	}

	shared.WriteJSON(w, http.StatusOK, doctor)
}

func (h *Handler) listDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.repo.List()
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "Error al listar doctores")
		return
	}

	shared.WriteJSON(w, http.StatusOK, doctors)
}

func (h *Handler) createDoctor(w http.ResponseWriter, r *http.Request) {
	var input CreateDoctorRequest

	if err := shared.ReadJSON(w, r, &input); err != nil {
		shared.WriteError(w, http.StatusBadRequest, "Json invalido")
		return
	}

	if err := input.Validar(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	doctor, err := h.repo.Create(input)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate key") {
			shared.WriteError(w, http.StatusConflict, "El correo ya existe")
			return
		}

		shared.WriteError(w, http.StatusInternalServerError, "Error al crear doctor")
		return
	}

	shared.WriteJSON(w, http.StatusCreated, doctor)
}
