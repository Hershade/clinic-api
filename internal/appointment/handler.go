package appointment

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

func (h *Handler) AppointmentsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAppointments(w, r)
	case http.MethodPost:
		h.createAppointment(w, r)
	default:
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
	}
}

func (h *Handler) AppointmentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/appointments/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	appointment, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteError(w, http.StatusNotFound, "cita no encontrada")
			return
		}
		log.Printf("error consultando cita por id: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al consultar cita")
		return
	}

	shared.WriteJSON(w, http.StatusOK, appointment)
}

func (h *Handler) listAppointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.repo.List()
	if err != nil {
		log.Printf("error listando citas: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al listar citas")
		return
	}

	shared.WriteJSON(w, http.StatusOK, appointments)
}

func (h *Handler) createAppointment(w http.ResponseWriter, r *http.Request) {
	var input CreateAppointmentRequest

	if err := shared.ReadJSON(w, r, &input); err != nil {
		log.Printf("error leyendo json de cita: %v", err)
		shared.WriteError(w, http.StatusBadRequest, "json invalido")
		return
	}

	if err := input.Validar(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	appointment, err := h.repo.Create(input)
	if err != nil {
		log.Printf("error creando cita: %v", err)

		switch {
		case errors.Is(err, ErrDoctorNotFound):
			shared.WriteError(w, http.StatusNotFound, "doctor no encontrado")
			return
		case errors.Is(err, ErrPatientNotFound):
			shared.WriteError(w, http.StatusNotFound, "paciente no encontrado")
			return
		case errors.Is(err, ErrAppointmentSlotTaken):
			shared.WriteError(w, http.StatusConflict, "el doctor ya tiene una cita en esa fecha y hora")
			return
		default:
			shared.WriteError(w, http.StatusInternalServerError, "error al crear cita")
			return
		}
	}

	shared.WriteJSON(w, http.StatusCreated, appointment)
}
