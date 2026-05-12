package main

import (
	"log"
	"net/http"
	"time"

	"clinic-api/internal/appointment"
	"clinic-api/internal/auth"
	"clinic-api/internal/doctor"
	"clinic-api/internal/health"
	"clinic-api/internal/middleware"
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

	authRepo := auth.NewRepository(database)
	if err := authRepo.EnsureSeedAdmin(cfg.AdminUsername, cfg.AdminPassword, cfg.AdminRole); err != nil {
		log.Fatalf("Error creando admin semilla: %v", err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresHours)
	authHandler := auth.NewHandler(authRepo, jwtManager)
	authMiddleware := auth.NewMiddleware(jwtManager)

	doctorRepo := doctor.NewRepository(database)
	doctorHandler := doctor.NewHandler(doctorRepo)

	patientRepo := patient.NewRepository(database)
	patientHandler := patient.NewHandler(patientRepo)

	appointmentRepo := appointment.NewRepository(database)
	appointmentHandler := appointment.NewHandler(appointmentRepo)

	mux := http.NewServeMux()

	// Públicas
	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("/auth/login", authHandler.Login)

	// Protegidas
	mux.Handle("/doctors", authMiddleware.RequireAuth(http.HandlerFunc(doctorHandler.DoctorsCollection)))
	mux.Handle("/doctors/", authMiddleware.RequireAuth(http.HandlerFunc(doctorHandler.DoctorByID)))

	mux.Handle("/patients", authMiddleware.RequireAuth(http.HandlerFunc(patientHandler.PatientsCollection)))
	mux.Handle("/patients/", authMiddleware.RequireAuth(http.HandlerFunc(patientHandler.PatientByID)))

	mux.Handle("/appointments", authMiddleware.RequireAuth(http.HandlerFunc(appointmentHandler.AppointmentsCollection)))
	mux.Handle("/appointments/", authMiddleware.RequireAuth(http.HandlerFunc(appointmentHandler.AppointmentRoutes)))

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Servidor corriendo en http://localhost:%s", cfg.AppPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
