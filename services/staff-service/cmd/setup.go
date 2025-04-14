package main

import (
	"log"

	"google.golang.org/grpc"

	coreDatabase "golang-microservices-boilerplate/pkg/core/database"
	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	"golang-microservices-boilerplate/services/staff-service/internal/controller"
	"golang-microservices-boilerplate/services/staff-service/internal/entity"
	staffRepoGorm "golang-microservices-boilerplate/services/staff-service/internal/repository"
	staffUseCase "golang-microservices-boilerplate/services/staff-service/internal/usecase"
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

	// Auto-migrate schema - Add lookup tables
	if err := db.MigrateModels(
		&entity.Staff{},
		&entity.ScheduleEntry{},
		&entity.Task{},
		&entity.StaffRole{},
		&entity.StaffStatus{},
		&entity.TaskStatus{},
	); err != nil {
		_ = db.Close()
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
	}
	logger.Info("Database schema migrated")
	return db
}

// setupDependencies initializes and returns the core dependencies: repositories, use case, mapper.
func setupDependencies(db *coreDatabase.DatabaseConnection, logger coreLogger.Logger) (staffUseCase.StaffUseCase, controller.Mapper) {
	// Instantiate all repositories
	staffRepo := staffRepoGorm.NewGormStaffRepository(db.DB)
	taskRepo := staffRepoGorm.NewGormTaskRepository(db.DB)
	staffRoleRepo := staffRepoGorm.NewGormStaffRoleRepository(db.DB)
	staffStatusRepo := staffRepoGorm.NewGormStaffStatusRepository(db.DB)
	taskStatusRepo := staffRepoGorm.NewGormTaskStatusRepository(db.DB)
	// Inject repositories into the use case
	uc := staffUseCase.NewStaffUseCase(staffRepo, taskRepo, staffRoleRepo, staffStatusRepo, taskStatusRepo, logger)
	mapper := controller.NewStaffMapper()
	return uc, mapper
}

// registerServices registers all gRPC services with the server.
func registerServices(s *grpc.Server, uc staffUseCase.StaffUseCase, mapper controller.Mapper) {
	controller.RegisterStaffServiceServer(s, uc, mapper) // Pass mapper
	// Register other services for this server if needed
}
