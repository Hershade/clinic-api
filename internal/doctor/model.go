package doctor

import (
	"errors"
	"net/mail"
	"strings"
)

// modelo de Doctor
type Doctor struct {
	ID           int64  `json:"id"`
	Nombre       string `json:"nombre"`
	Especialidad string `json:"specialidad"`
	Telefono     string `json:"telefono"`
	Correo       string `json:"correo"`
	Activo       bool   `json:"activo"`
}

// modelo del request
type CreateDoctorRequest struct {
	Nombre       string `json:"nombre"`
	Especialidad string `json:"especialidad"`
	Telefono     string `json:"telefono"`
	Correo       string `json:"correo"`
}

func (r *CreateDoctorRequest) Validar() error {
	r.Nombre = strings.TrimSpace(r.Nombre)
	r.Especialidad = strings.TrimSpace(r.Especialidad)
	r.Telefono = strings.TrimSpace(r.Telefono)
	r.Correo = strings.TrimSpace(r.Correo)

	if r.Nombre == "" {
		return errors.New("El nombre es obligatorio")
	}

	if r.Especialidad == "" {
		return errors.New("La especialidad es obligatoria")
	}

	if r.Telefono == "" {
		return errors.New("El telefono es obligatorio")
	}

	if r.Correo == "" {
		return errors.New("El correo es obligatorio")
	}

	if len(r.Nombre) > 100 {
		return errors.New("El nombre es demasiado largo")
	}

	if len(r.Especialidad) > 100 {
		return errors.New("La especialidad es demasiado larga")
	}

	if _, err := mail.ParseAddress(r.Correo); err != nil {
		return errors.New("El correo no es valido")
	}

	return nil
}
