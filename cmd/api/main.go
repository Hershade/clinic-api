package main

import (
	"log"
	"net/http"
	"time"

	"clinic-api/internal/appointment"
	"clinic-api/internal/doctor"
	"clinic-api/internal/health"
	"clinic-api/internal/patient"
	"clinic-api/internal/platform/config"
	"clinic-api/internal/platform/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuracion: %v", err)
	}

	database, err := db.OpenPostgres(cfg)
	if err != nil {
		log.Fatalf("Error conectando a postgres: %v", err)
	}

	defer database.Close()

	healthHandler := health.NewHandler(database)

	doctorRepo := doctor.NewRepository(database)
	doctorHandler := doctor.NewHandler(doctorRepo)

	patientRepo := patient.NewRepository(database)
	patientHandler := patient.NewHandler(patientRepo)

	appointmentRepo := appointment.NewRepository(database)
	appointmentHandler := appointment.NewHandler(appointmentRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.Check)

	mux.HandleFunc("/doctors", doctorHandler.DoctorsCollection)
	mux.HandleFunc("/doctors/", doctorHandler.DoctorByID)

	mux.HandleFunc("/patients", patientHandler.PatientsCollection)
	mux.HandleFunc("/patients/", patientHandler.PatientByID)

	mux.HandleFunc("/appointments", appointmentHandler.AppointmentsCollection)
	mux.HandleFunc("/appointments/", appointmentHandler.AppointmentRoutes)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Servidor corriendo en http://localhost:%s", cfg.AppPort)

	if err := server.ListenAndServe(); err != nil {

		log.Fatal(err)
	}
}
