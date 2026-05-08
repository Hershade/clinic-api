package patient

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"clinic-api/internal/shared"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) PatientsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listPatients(w, r)
	case http.MethodPost:
		h.createPatient(w, r)
	default:
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
	}
}

func (h *Handler) PatientByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/patients/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	patient, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteError(w, http.StatusNotFound, "paciente no encontrado")
			return
		}
		shared.WriteError(w, http.StatusInternalServerError, "error al consultar paciente")
		return
	}

	shared.WriteJSON(w, http.StatusOK, patient)
}

func (h *Handler) listPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.repo.List()
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, "error al listar pacientes")
		return
	}

	shared.WriteJSON(w, http.StatusOK, patients)
}

func (h *Handler) createPatient(w http.ResponseWriter, r *http.Request) {
	var input CreatePatientRequest

	if err := shared.ReadJSON(w, r, &input); err != nil {
		log.Printf("error leyendo json de paciente: %v", err)
		shared.WriteError(w, http.StatusBadRequest, "json invalido")
		return
	}

	if err := input.Validar(); err != nil {
		log.Printf("error validadndo paciente: %v", err)
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	patient, err := h.repo.Create(input)
	if err != nil {
		log.Printf("error creando paciente: %v", err)

		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			shared.WriteError(w, http.StatusConflict, "el dpi o correo ya existe")
			return
		}

		shared.WriteError(w, http.StatusInternalServerError, "error al crear paciente")
		return
	}

	shared.WriteJSON(w, http.StatusCreated, patient)
}
