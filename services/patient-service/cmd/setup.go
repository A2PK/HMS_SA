package main

import (
	"log"

	"google.golang.org/grpc"

	coreDatabase "golang-microservices-boilerplate/pkg/core/database"
	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	"golang-microservices-boilerplate/services/patient-service/internal/controller"
	"golang-microservices-boilerplate/services/patient-service/internal/entity"
	patientRepoGorm "golang-microservices-boilerplate/services/patient-service/internal/repository"
	patientUseCase "golang-microservices-boilerplate/services/patient-service/internal/usecase"
)

// setupLogger initializes the logger based on environment configuration.
func setupLogger(appName string) coreLogger.Logger {
	logConfig := coreLogger.LoadLogConfigFromEnv()
	logConfig.AppName = appName
	logger, err := coreLogger.NewLogger(logConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Info("Logger initialized", "level", logConfig.Level, "format", logConfig.Format)
	return logger
}

// setupDatabase initializes the database connection and performs migrations.
func setupDatabase(logger coreLogger.Logger) *coreDatabase.DatabaseConnection {
	dbConfig := coreDatabase.DefaultDBConfig() // Load config from env or defaults
	db, err := coreDatabase.NewDatabaseConnection(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	logger.Info("Database connection established")

	// Auto-migrate schema
	if err := db.MigrateModels(&entity.Patient{}, &entity.MedicalRecord{}); err != nil {
		// Attempt close before fatal
		_ = db.Close()
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
	}
	logger.Info("Database schema migrated")
	return db
}

// setupDependencies initializes and returns the core dependencies: repository, use case, mapper.
func setupDependencies(db *coreDatabase.DatabaseConnection, logger coreLogger.Logger) (patientUseCase.PatientUseCase, controller.Mapper) {
	repo := patientRepoGorm.NewGormPatientRepository(db.DB) // Assuming Gorm implementation
	uc := patientUseCase.NewPatientUseCase(repo, logger)
	mapper := controller.NewPatientMapper()
	return uc, mapper
}

// registerServices registers all gRPC services with the server.
func registerServices(s *grpc.Server, uc patientUseCase.PatientUseCase, mapper controller.Mapper) {
	controller.RegisterPatientServiceServer(s, uc, mapper)
	// Register other services for this server if needed
}
