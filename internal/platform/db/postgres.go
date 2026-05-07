package db

import (
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/platform/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// lectura de variables de entorno
func OpenPostgres(cfg *config.EnvConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)                  //maximo de conexiones abiertas
	db.SetMaxIdleConns(5)                   // conexiones inactivas que se pueden reutilizar
	db.SetConnMaxLifetime(30 * time.Minute) //tiempo maximo de vida de una conexion

	//hacer ping para validar la conexion
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
