package main

import (
	"log"

	"google.golang.org/grpc"

	coreDatabase "golang-microservices-boilerplate/pkg/core/database"
	coreGrpc "golang-microservices-boilerplate/pkg/core/grpc"
	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	staff_pb "golang-microservices-boilerplate/proto/staff-service"
	"golang-microservices-boilerplate/services/appointment-service/internal/controller"
	"golang-microservices-boilerplate/services/appointment-service/internal/entity"
	appointmentRepoGorm "golang-microservices-boilerplate/services/appointment-service/internal/repository"
	appointmentUseCase "golang-microservices-boilerplate/services/appointment-service/internal/usecase"
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
	dbConfig := coreDatabase.DefaultDBConfig()
	db, err := coreDatabase.NewDatabaseConnection(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	logger.Info("Database connection established")

	// Auto-migrate schema
	if err := db.MigrateModels(&entity.Appointment{}); err != nil {
		_ = db.Close()
		logger.Fatal("Failed to auto-migrate database schema", "error", err)
	}
	logger.Info("Database schema migrated")
	return db
}

// setupDependencies initializes and returns the core dependencies: repository, use case, mapper, and staff client adapter.
func setupDependencies(db *coreDatabase.DatabaseConnection, logger coreLogger.Logger, staffServiceAddress string) (appointmentUseCase.AppointmentUseCase, controller.Mapper) {
	// Staff Service gRPC Client
	staffClientConn, err := coreGrpc.NewBaseGrpcClient(logger, &coreGrpc.GrpcClientConfig{
		ServiceHost:            staffServiceAddress,
		AllowInsecureTransport: true,
	})
	if err != nil {
		logger.Fatal("Failed to connect to staff service", "address", staffServiceAddress, "error", err)
	}

	staffServiceClient := staff_pb.NewStaffServiceClient(staffClientConn.Conn)

	// Create adapter
	staffAdapter := appointmentUseCase.NewGrpcStaffServiceClientAdapter(staffServiceClient, logger)

	// Init other dependencies
	repo := appointmentRepoGorm.NewGormAppointmentRepository(db.DB)
	uc := appointmentUseCase.NewAppointmentUseCase(repo, staffAdapter, logger)
	mapper := controller.NewAppointmentMapper() // Instantiate the correct mapper

	return uc, mapper
}

// registerServices registers all gRPC services with the server.
func registerServices(s *grpc.Server, uc appointmentUseCase.AppointmentUseCase, mapper controller.Mapper) {
	controller.RegisterAppointmentServiceServer(s, uc, mapper)
}
