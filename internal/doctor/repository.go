package doctor

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List() ([]Doctor, error) {
	rows, err := r.db.Query(`
		SELECT id, nombre, especialidad, telefono, correo, activo
		FROM doctors
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := []Doctor{}

	for rows.Next() {
		var d Doctor
		err := rows.Scan(
			&d.ID,
			&d.Nombre,
			&d.Especialidad,
			&d.Telefono,
			&d.Correo,
			&d.Activo,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}

	return doctors, rows.Err()
}

func (r *Repository) GetByID(id int64) (Doctor, error) {
	var d Doctor

	err := r.db.QueryRow(`
		SELECT id, nombre, especialidad, telefono, correo,activo
		FROM doctors
		WHERE id = $1
	`, id).Scan(
		&d.ID,
		&d.Nombre,
		&d.Especialidad,
		&d.Telefono,
		&d.Correo,
		&d.Activo,
	)

	return d, err
}

func (r *Repository) Create(input CreateDoctorRequest) (Doctor, error) {
	var d Doctor

	err := r.db.QueryRow(`
		INSERT INTO doctors (nombre, especialidad, telefono, correo)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nombre, especialidad, telefono, correo, activo
	`,
		input.Nombre,
		input.Especialidad,
		input.Telefono,
		input.Correo,
	).Scan(
		&d.ID,
		&d.Nombre,
		&d.Especialidad,
		&d.Telefono,
		&d.Correo,
		&d.Activo,
	)

	return d, err
}
