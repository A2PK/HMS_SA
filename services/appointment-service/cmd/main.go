package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// Core packages

	coreGrpc "golang-microservices-boilerplate/pkg/core/grpc"
	"golang-microservices-boilerplate/pkg/utils"
)

func main() {
	// --- Configuration & Basic Setup ---
	if err := utils.LoadEnv(); err != nil {
		log.Printf("Warning: .env file not found or error loading, using environment variables: %v", err)
	}
	appName := utils.GetEnv("SERVER_APP_NAME", "Appointment Service")
	staffServiceAddress := utils.GetEnv("STAFF_SERVICE_ADDRESS", "localhost:50052") // Example address

	// --- Setup Dependencies using functions from setup.go ---
	logger := setupLogger(appName)
	logger.Info("Appointment service starting...")

	db := setupDatabase(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Error closing database connection", "error", err)
		}
	}()

	// Note: setupDependencies now handles staff client creation
	uc, mapper := setupDependencies(db, logger, staffServiceAddress)

	// --- Setup gRPC Server (using coreGrpc helper) ---
	grpcServer := coreGrpc.NewBaseGrpcServer(logger)

	// Register Services using function from setup.go
	registerServices(grpcServer.Server(), uc, mapper) // Pass mapper
	logger.Info("Appointment gRPC service registered")

	// Health check and Reflection are typically handled within NewBaseGrpcServer or its Start method

	// --- Start Server ---
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", "error", err)
	}
	logger.Info("gRPC server started successfully", "address", grpcServer.Config.Host+":"+grpcServer.Config.Port)

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	grpcServer.Stop()                         // Stop the gRPC server
	logger.Info("Server stopped gracefully.") // db closed by defer

	logger.Info("Appointment service stopped.")
}
