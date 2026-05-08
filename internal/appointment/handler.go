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

func (h *Handler) AppointmentRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/appointments/")
	path = strings.Trim(path, "/")

	if path == "" {
		shared.WriteError(w, http.StatusNotFound, "ruta no encontrada")
		return
	}

	parts := strings.Split(path, "/")

	switch {
	case len(parts) == 1:
		h.getAppointmentByID(w, r, parts[0])
	case len(parts) == 2 && parts[0] == "doctor":
		h.listAppointmentsByDoctor(w, r, parts[1])
	case len(parts) == 2 && parts[0] == "patient":
		h.listAppointmentsByPatient(w, r, parts[1])
	case len(parts) == 2 && parts[1] == "cancel":
		h.cancelAppointment(w, r, parts[0])
	default:
		shared.WriteError(w, http.StatusNotFound, "ruta no encontrada")
	}
}

func (h *Handler) getAppointmentByID(w http.ResponseWriter, r *http.Request, idStr string) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

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

func (h *Handler) listAppointmentsByDoctor(w http.ResponseWriter, r *http.Request, idStr string) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	doctorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || doctorID <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "doctor id invalido")
		return
	}

	appointments, err := h.repo.ListByDoctor(doctorID)
	if err != nil {
		log.Printf("error listando citas por doctor: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al listar citas por doctor")
		return
	}

	shared.WriteJSON(w, http.StatusOK, appointments)
}

func (h *Handler) listAppointmentsByPatient(w http.ResponseWriter, r *http.Request, idStr string) {
	if r.Method != http.MethodGet {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	patientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || patientID <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "patient id invalido")
		return
	}

	appointments, err := h.repo.ListByPatient(patientID)
	if err != nil {
		log.Printf("error listando citas por paciente: %v", err)
		shared.WriteError(w, http.StatusInternalServerError, "error al listar citas por paciente")
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

func (h *Handler) cancelAppointment(w http.ResponseWriter, r *http.Request, idStr string) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPost {
		shared.WriteError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		shared.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	appointment, err := h.repo.Cancel(id)
	if err != nil {
		log.Printf("error cancelando cita: %v", err)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			shared.WriteError(w, http.StatusNotFound, "cita no encontrada")
			return
		case errors.Is(err, ErrAppointmentAlreadyCanceled):
			shared.WriteError(w, http.StatusConflict, "la cita ya fue cancelada")
			return
		default:
			shared.WriteError(w, http.StatusInternalServerError, "error al cancelar cita")
			return
		}
	}

	shared.WriteJSON(w, http.StatusOK, appointment)
}
