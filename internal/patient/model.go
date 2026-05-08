package patient

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

type Patient struct {
	ID              int64  `json:"id"`
	Nombre          string `json:"nombre"`
	DPI             string `json:"dpi"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Activo          bool   `json:"activo"`
}

type CreatePatientRequest struct {
	Nombre          string `json:"nombre"`
	DPI             string `json:"dpi"`
	Telefono        string `json:"telefono"`
	Correo          string `json:"correo"`
	FechaNacimiento string `json:"fecha_nacimiento"`
}

func (r *CreatePatientRequest) Validar() error {
	r.Nombre = strings.TrimSpace(r.Nombre)
	r.DPI = strings.TrimSpace(r.DPI)
	r.Telefono = strings.TrimSpace(r.Telefono)
	r.Correo = strings.TrimSpace(r.Correo)
	r.FechaNacimiento = strings.TrimSpace(r.FechaNacimiento)

	if r.Nombre == "" {
		return errors.New("el nombre es obligatorio")
	}

	if r.DPI == "" {
		return errors.New("el dpi es obligatorio")
	}

	if r.Telefono == "" {
		return errors.New("el telefono es obligatorio")
	}

	if r.Correo == "" {
		return errors.New("el correo es obligatorio")
	}

	if r.FechaNacimiento == "" {
		return errors.New("la fecha de nacimineto es obligatorio")
	}

	if len(r.Nombre) > 100 {
		return errors.New("el nombre es demasido largo")
	}

	if len(r.DPI) > 20 {
		return errors.New("el dpi es demasiado largo")
	}

	if _, err := mail.ParseAddress(r.Correo); err != nil {
		return errors.New("el correo no es valido")
	}

	if _, err := time.Parse("2006-01-02", r.FechaNacimiento); err != nil {
		return errors.New("la fecha de nacimiento debe tener formato YYYY-MM-DD")
	}

	return nil
}
