package appointment

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

func (r *Repository) List() ([]Appointment, error) {
	return r.listByQuery(`
		SELECT
			a.id,
			a.doctor_id,
			d.nombre AS doctor_nombre,
			a.patient_id,
			p.nombre AS patient_nombre,
			TO_CHAR(a.fecha, 'YYYY-MM-DD') AS fecha,
			TO_CHAR(a.hora, 'HH24:MI:SS') AS hora,
			a.motivo,
			a.estado
		FROM appointments a
		INNER JOIN doctors d ON d.id = a.doctor_id
		INNER JOIN patients p ON p.id = a.patient_id
		ORDER BY a.fecha ASC, a.hora ASC, a.id ASC
	`)
}

func (r *Repository) ListByDoctor(doctorID int64) ([]Appointment, error) {
	return r.listByQuery(`
		SELECT
			a.id,
			a.doctor_id,
			d.nombre AS doctor_nombre,
			a.patient_id,
			p.nombre AS patient_nombre,
			TO_CHAR(a.fecha, 'YYYY-MM-DD') AS fecha,
			TO_CHAR(a.hora, 'HH24:MI:SS') AS hora,
			a.motivo,
			a.estado
		FROM appointments a
		INNER JOIN doctors d ON d.id = a.doctor_id
		INNER JOIN patients p ON p.id = a.patient_id
		WHERE a.doctor_id = $1
		ORDER BY a.fecha ASC, a.hora ASC, a.id ASC
	`, doctorID)
}

func (r *Repository) ListByPatient(patientID int64) ([]Appointment, error) {
	return r.listByQuery(`
		SELECT
			a.id,
			a.doctor_id,
			d.nombre AS doctor_nombre,
			a.patient_id,
			p.nombre AS patient_nombre,
			TO_CHAR(a.fecha, 'YYYY-MM-DD') AS fecha,
			TO_CHAR(a.hora, 'HH24:MI:SS') AS hora,
			a.motivo,
			a.estado
		FROM appointments a
		INNER JOIN doctors d ON d.id = a.doctor_id
		INNER JOIN patients p ON p.id = a.patient_id
		WHERE a.patient_id = $1
		ORDER BY a.fecha ASC, a.hora ASC, a.id ASC
	`, patientID)
}

func (r *Repository) listByQuery(query string, args ...any) ([]Appointment, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error en query de appointments: %w", err)
	}
	defer rows.Close()

	appointments := []Appointment{}

	for rows.Next() {
		var a Appointment

		err := rows.Scan(
			&a.ID,
			&a.DoctorID,
			&a.DoctorNombre,
			&a.PatientID,
			&a.PatientNombre,
			&a.Fecha,
			&a.Hora,
			&a.Motivo,
			&a.Estado,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando appointment: %w", err)
		}

		appointments = append(appointments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando rows de appointments: %w", err)
	}

	return appointments, nil
}

func (r *Repository) GetByID(id int64) (Appointment, error) {
	var a Appointment

	err := r.db.QueryRow(`
		SELECT
			a.id,
			a.doctor_id,
			d.nombre AS doctor_nombre,
			a.patient_id,
			p.nombre AS patient_nombre,
			TO_CHAR(a.fecha, 'YYYY-MM-DD') AS fecha,
			TO_CHAR(a.hora, 'HH24:MI:SS') AS hora,
			a.motivo,
			a.estado
		FROM appointments a
		INNER JOIN doctors d ON d.id = a.doctor_id
		INNER JOIN patients p ON p.id = a.patient_id
		WHERE a.id = $1
	`, id).Scan(
		&a.ID,
		&a.DoctorID,
		&a.DoctorNombre,
		&a.PatientID,
		&a.PatientNombre,
		&a.Fecha,
		&a.Hora,
		&a.Motivo,
		&a.Estado,
	)

	return a, err
}

func (r *Repository) Create(input CreateAppointmentRequest) (Appointment, error) {
	doctorExists, err := r.doctorExists(input.DoctorID)
	if err != nil {
		return Appointment{}, err
	}
	if !doctorExists {
		return Appointment{}, ErrDoctorNotFound
	}

	patientExists, err := r.patientExists(input.PatientID)
	if err != nil {
		return Appointment{}, err
	}
	if !patientExists {
		return Appointment{}, ErrPatientNotFound
	}

	taken, err := r.slotTaken(input.DoctorID, input.Fecha, input.Hora)
	if err != nil {
		return Appointment{}, err
	}
	if taken {
		return Appointment{}, ErrAppointmentSlotTaken
	}

	var id int64

	err = r.db.QueryRow(`
		INSERT INTO appointments (doctor_id, patient_id, fecha, hora, motivo)
		VALUES ($1, $2, $3::date, $4::time, $5)
		RETURNING id
	`,
		input.DoctorID,
		input.PatientID,
		input.Fecha,
		input.Hora,
		input.Motivo,
	).Scan(&id)
	if err != nil {
		return Appointment{}, fmt.Errorf("error creando appointment: %w", err)
	}

	return r.GetByID(id)
}

func (r *Repository) Cancel(id int64) (Appointment, error) {
	current, err := r.GetByID(id)
	if err != nil {
		return Appointment{}, err
	}

	if current.Estado == "cancelada" {
		return Appointment{}, ErrAppointmentAlreadyCanceled
	}

	_, err = r.db.Exec(`
		UPDATE appointments
		SET estado = 'cancelada'
		WHERE id = $1
	`, id)
	if err != nil {
		return Appointment{}, fmt.Errorf("error cancelando appointment: %w", err)
	}

	return r.GetByID(id)
}

func (r *Repository) doctorExists(id int64) (bool, error) {
	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM doctors
			WHERE id = $1 AND activo = TRUE
		)
	`, id).Scan(&exists)

	return exists, err
}

func (r *Repository) patientExists(id int64) (bool, error) {
	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM patients
			WHERE id = $1 AND activo = TRUE
		)
	`, id).Scan(&exists)

	return exists, err
}

func (r *Repository) slotTaken(doctorID int64, fecha, hora string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM appointments
			WHERE doctor_id = $1
			  AND fecha = $2::date
			  AND hora = $3::time
			  AND estado IN ('pendiente', 'confirmada')
		)
	`, doctorID, fecha, hora).Scan(&exists)

	return exists, err
}
