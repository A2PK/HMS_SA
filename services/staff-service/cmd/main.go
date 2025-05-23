package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// Core packages

	coreGrpc "golang-microservices-boilerplate/pkg/core/grpc"
	"golang-microservices-boilerplate/pkg/utils"
	// Staff service internal packages
	// For AutoMigrate
	// Assuming GORM implementation
)

func main() {
	// --- Configuration & Basic Setup ---
	if err := utils.LoadEnv(); err != nil {
		log.Printf("Warning: .env file not found or error loading, using environment variables: %v", err)
	}
	appName := utils.GetEnv("SERVER_APP_NAME", "Staff Service")

	// --- Setup Dependencies using functions from setup.go ---
	logger := setupLogger(appName)
	logger.Info("Staff service starting...")

	db := setupDatabase(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Error closing database connection", "error", err)
		}
	}()

	uc, mapper := setupDependencies(db, logger)

	// --- Setup gRPC Server (using coreGrpc helper) ---
	grpcServer := coreGrpc.NewBaseGrpcServer(logger)

	// Register Services using function from setup.go
	registerServices(grpcServer.Server(), uc, mapper)
	logger.Info("Staff gRPC service registered")

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

	logger.Info("Staff service stopped.")
}
