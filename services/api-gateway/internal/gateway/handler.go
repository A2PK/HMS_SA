package gateway

import (
	"fmt"

	appointment_pb "golang-microservices-boilerplate/proto/appointment-service"
	patient_pb "golang-microservices-boilerplate/proto/patient-service"
	staff_pb "golang-microservices-boilerplate/proto/staff-service"
	user_pb "golang-microservices-boilerplate/proto/user-service"
	"golang-microservices-boilerplate/services/api-gateway/internal/domain"
)

// setupUserServiceHandlers registers handlers for the user service
func (g *Gateway) setupUserServiceHandlers(service domain.Service) error {
	err := user_pb.RegisterUserServiceHandlerFromEndpoint(g.ctx, g.gwMux, service.Endpoint, g.opts)
	if err != nil {
		g.logger.Error("Failed to register user service handler from endpoint", "endpoint", service.Endpoint, "error", err)
		return fmt.Errorf("failed to register user service handler from endpoint %s: %w", service.Endpoint, err)
	}

	g.logger.Info("Registered gRPC-Gateway handlers via endpoint", "service", "user-service", "endpoint", service.Endpoint)
	return nil
}

func (g *Gateway) setupPatientServiceHandlers(service domain.Service) error {
	err := patient_pb.RegisterPatientServiceHandlerFromEndpoint(g.ctx, g.gwMux, service.Endpoint, g.opts)
	if err != nil {
		g.logger.Error("Failed to register patient service handler from endpoint", "endpoint", service.Endpoint, "error", err)
		return fmt.Errorf("failed to register patient service handler from endpoint %s: %w", service.Endpoint, err)
	}
	g.logger.Info("Registered gRPC-Gateway handlers via endpoint", "service", "patient-service", "endpoint", service.Endpoint)
	return nil
}

func (g *Gateway) setupAppointmentServiceHandlers(service domain.Service) error {
	err := appointment_pb.RegisterAppointmentServiceHandlerFromEndpoint(g.ctx, g.gwMux, service.Endpoint, g.opts)
	if err != nil {
		g.logger.Error("Failed to register appointment service handler from endpoint", "endpoint", service.Endpoint, "error", err)
		return fmt.Errorf("failed to register appointment service handler from endpoint %s: %w", service.Endpoint, err)
	}
	g.logger.Info("Registered gRPC-Gateway handlers via endpoint", "service", "appointment-service", "endpoint", service.Endpoint)
	return nil
}

func (g *Gateway) setupStaffServiceHandlers(service domain.Service) error {
	err := staff_pb.RegisterStaffServiceHandlerFromEndpoint(g.ctx, g.gwMux, service.Endpoint, g.opts)
	if err != nil {
		g.logger.Error("Failed to register staff service handler from endpoint", "endpoint", service.Endpoint, "error", err)
		return fmt.Errorf("failed to register staff service handler from endpoint %s: %w", service.Endpoint, err)
	}
	g.logger.Info("Registered gRPC-Gateway handlers via endpoint", "service", "staff-service", "endpoint", service.Endpoint)
	return nil
}
