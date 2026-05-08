package patient

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List() ([]Patient, error) {
	rows, err := r.db.Query(`
		SELECT id, nombre, dpi, telefono, correo,
		       TO_CHAR(fecha_nacimiento, 'YYYY-MM-DD') AS fecha_nacimiento,
		       activo
		FROM patients
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("error en query de patients: %w", err)
	}
	defer rows.Close()

	patients := []Patient{}

	for rows.Next() {
		var p Patient

		err := rows.Scan(
			&p.ID,
			&p.Nombre,
			&p.DPI,
			&p.Telefono,
			&p.Correo,
			&p.FechaNacimiento,
			&p.Activo,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando patient: %w", err)
		}

		patients = append(patients, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando rows de patients: %w", err)
	}

	return patients, nil
}

func (r *Repository) GetByID(id int64) (Patient, error) {
	var p Patient

	err := r.db.QueryRow(`
		SELECT id, nombre, dpi, telefono, correo,
		       TO_CHAR(fecha_nacimiento, 'YYYY-MM-DD') as fecha_nacimiento,
		       activo
		FROM patients
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.Nombre,
		&p.DPI,
		&p.Telefono,
		&p.Correo,
		&p.FechaNacimiento,
		&p.Activo,
	)

	return p, err
}

func (r *Repository) Create(input CreatePatientRequest) (Patient, error) {
	var p Patient

	err := r.db.QueryRow(`
		INSERT INTO patients (nombre, dpi, telefono, correo, fecha_nacimiento)
		VALUES ($1, $2, $3, $4, $5::date)
		RETURNING id, nombre, dpi, telefono, correo,
			TO_CHAR(fecha_nacimiento, 'YYYY-MM-DD') AS fecha_nacimiento,
		activo
	`,
		input.Nombre,
		input.DPI,
		input.Telefono,
		input.Correo,
		input.FechaNacimiento,
	).Scan(
		&p.ID,
		&p.Nombre,
		&p.DPI,
		&p.Telefono,
		&p.Correo,
		&p.FechaNacimiento,
		&p.Activo,
	)

	if err != nil {
		return p, fmt.Errorf("error en create patient: %w", err)
	}

	return p, err
}
