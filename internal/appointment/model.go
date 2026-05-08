package appointment

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrDoctorNotFound       = errors.New("doctor no encontrado")
	ErrPatientNotFound      = errors.New("paciente no encontrado")
	ErrAppointmentSlotTaken = errors.New("el doctor ya tiene una cita en esa fecha y hora")
)

type Appointment struct {
	ID            int64  `json:"id"`
	DoctorID      int64  `json:"doctor_id"`
	DoctorNombre  string `json:"doctor_nombre"`
	PatientID     int64  `json:"patient_id"`
	PatientNombre string `json:"patient_nombre"`
	Fecha         string `json:"fecha"`
	Hora          string `json:"hora"`
	Motivo        string `json:"motivo"`
	Estado        string `json:"estado"`
}

type CreateAppointmentRequest struct {
	DoctorID  int64  `json:"doctor_id"`
	PatientID int64  `json:"patient_id"`
	Fecha     string `json:"fecha"`
	Hora      string `json:"hora"`
	Motivo    string `json:"motivo"`
}

func (r *CreateAppointmentRequest) Validar() error {
	r.Fecha = strings.TrimSpace(r.Fecha)
	r.Hora = strings.TrimSpace(r.Hora)
	r.Motivo = strings.TrimSpace(r.Motivo)

	if r.DoctorID <= 0 {
		return errors.New("doctor_id es obligatorio")
	}

	if r.PatientID <= 0 {
		return errors.New("patient_id es obligatorio")
	}

	if r.Fecha == "" {
		return errors.New("la fecha es obligatoria")
	}

	if r.Hora == "" {
		return errors.New("la hora es obligatoria")
	}

	if r.Motivo == "" {
		return errors.New("el motivo es obligatorio")
	}

	if len(r.Motivo) > 255 {
		return errors.New("el motivo es demasiado largo")
	}

	if _, err := time.Parse("2006-01-02", r.Fecha); err != nil {
		return errors.New("la fecha debe tener formato YYYY-MM-DD")
	}

	horaNormalizada, err := normalizarHora(r.Hora)
	if err != nil {
		return errors.New("la hora debe tener formato HH:MM o HH:MM:SS")
	}
	r.Hora = horaNormalizada

	return nil
}

func normalizarHora(h string) (string, error) {
	if t, err := time.Parse("15:04", h); err == nil {
		return t.Format("15:04:05"), nil
	}

	if t, err := time.Parse("15:04:05", h); err == nil {
		return t.Format("15:04:05"), nil
	}

	return "", errors.New("hora invalida")
}
