package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) EnsureSeedAdmin(username, password, role string) error {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1
		    FROM users
		    WHERE username = $1
		)
	`, username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error verificando usuario admin: %w", err)
	}

	if exists {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generando hash del admin: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO users (username, password_hash, role, activo)
		VALUES ($1, $2, $3, TRUE)
	`, username, string(hash), role)
	if err != nil {
		return fmt.Errorf("error creando usuario admin semilla: %w", err)
	}

	return nil
}

func (r *Repository) FindByUsername(username string) (User, error) {
	var user User

	err := r.db.QueryRow(`
		SELECT id, username, password_hash, role, activo
		FROM users
		WHERE username = $1
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.Activo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, sql.ErrNoRows
		}
		return User{}, fmt.Errorf("error consultando usuario por username: %w", err)
	}
	return user, nil
}
