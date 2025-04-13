package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	core_usecase "golang-microservices-boilerplate/pkg/core/usecase"
	pb "golang-microservices-boilerplate/proto/appointment-service" // Used indirectly
	"golang-microservices-boilerplate/services/appointment-service/internal/usecase"
)

// AppointmentServer defines the interface for the gRPC service handler.
type AppointmentServer interface {
	pb.AppointmentServiceServer // Embed generated interface
}

// Ensure appointmentServer implements AppointmentServer.
var _ AppointmentServer = (*appointmentServer)(nil)

// --- gRPC Server Implementation ---

type appointmentServer struct {
	pb.UnimplementedAppointmentServiceServer // Embed for forward compatibility
	uc                                       usecase.AppointmentUseCase
	mapper                                   Mapper // Use interface
}

// NewAppointmentServer creates a new gRPC server instance.
// Accepts Mapper interface and returns AppointmentServer interface.
func NewAppointmentServer(uc usecase.AppointmentUseCase, mapper Mapper) AppointmentServer {
	return &appointmentServer{
		uc:     uc,
		mapper: mapper,
	}
}

// RegisterAppointmentServiceServer registers the appointment service implementation with the gRPC server.
// Accepts use case and mapper.
func RegisterAppointmentServiceServer(s *grpc.Server, uc usecase.AppointmentUseCase, mapper Mapper) {
	server := NewAppointmentServer(uc, mapper)
	pb.RegisterAppointmentServiceServer(s, server)
}

// ScheduleAppointment implements the corresponding gRPC method.
func (s *appointmentServer) ScheduleAppointment(ctx context.Context, req *pb.ScheduleAppointmentRequest) (*pb.ScheduleAppointmentResponse, error) {
	// Use case already accepts the proto request, no mapping needed here.
	aptEntity, err := s.uc.ScheduleAppointment(ctx, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProto, err := s.mapper.EntityToProto(aptEntity)
	if err != nil {
		// Log internal mapping error
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.ScheduleAppointmentResponse{Appointment: aptProto}, nil
}

// GetAppointmentDetails implements the corresponding gRPC method.
func (s *appointmentServer) GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error) {
	appointmentID, err := uuid.Parse(req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid appointment ID format: %v", err)
	}

	aptEntity, err := s.uc.GetAppointmentDetails(ctx, appointmentID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProto, err := s.mapper.EntityToProto(aptEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.GetAppointmentDetailsResponse{Appointment: aptProto}, nil
}

// UpdateAppointmentStatus implements the corresponding gRPC method.
func (s *appointmentServer) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error) {
	appointmentID, err := uuid.Parse(req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid appointment ID format: %v", err)
	}

	// Use case already accepts the proto request.
	aptEntity, err := s.uc.UpdateAppointmentStatus(ctx, appointmentID, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProto, err := s.mapper.EntityToProto(aptEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.UpdateAppointmentStatusResponse{Appointment: aptProto}, nil
}

// RescheduleAppointment implements the corresponding gRPC method.
func (s *appointmentServer) RescheduleAppointment(ctx context.Context, req *pb.RescheduleAppointmentRequest) (*pb.RescheduleAppointmentResponse, error) {
	appointmentID, err := uuid.Parse(req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid appointment ID format: %v", err)
	}

	// Use case already accepts the proto request.
	aptEntity, err := s.uc.RescheduleAppointment(ctx, appointmentID, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProto, err := s.mapper.EntityToProto(aptEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.RescheduleAppointmentResponse{Appointment: aptProto}, nil
}

// CancelAppointment implements the corresponding gRPC method.
func (s *appointmentServer) CancelAppointment(ctx context.Context, req *pb.CancelAppointmentRequest) (*emptypb.Empty, error) {
	appointmentID, err := uuid.Parse(req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid appointment ID format: %v", err)
	}

	err = s.uc.CancelAppointment(ctx, appointmentID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

// GetAppointmentsForPatient implements the corresponding gRPC method.
func (s *appointmentServer) GetAppointmentsForPatient(ctx context.Context, req *pb.GetAppointmentsForPatientRequest) (*pb.GetAppointmentsForPatientResponse, error) {
	patientID, err := uuid.Parse(req.PatientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid patient ID format: %v", err)
	}

	aptEntities, err := s.uc.GetAppointmentsForPatient(ctx, patientID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProtos, err := s.mapper.EntitiesToProto(aptEntities)
	if err != nil {
		// This indicates an internal issue with mapping potentially valid data
		return nil, status.Errorf(codes.Internal, "failed to map results to proto: %v", err)
	}

	return &pb.GetAppointmentsForPatientResponse{Appointments: aptProtos}, nil
}

// GetAppointmentsForDoctor implements the corresponding gRPC method.
func (s *appointmentServer) GetAppointmentsForDoctor(ctx context.Context, req *pb.GetAppointmentsForDoctorRequest) (*pb.GetAppointmentsForDoctorResponse, error) {
	doctorID, err := uuid.Parse(req.DoctorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid doctor ID format: %v", err)
	}
	startTime := time.Time{}
	if req.StartTime != nil {
		startTime = req.StartTime.AsTime()
	}
	endTime := time.Time{}
	if req.EndTime != nil {
		endTime = req.EndTime.AsTime()
	}

	if startTime.IsZero() || endTime.IsZero() || endTime.Before(startTime) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time range provided")
	}

	aptEntities, err := s.uc.GetAppointmentsForDoctor(ctx, doctorID, startTime, endTime)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	aptProtos, err := s.mapper.EntitiesToProto(aptEntities)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map results to proto: %v", err)
	}

	return &pb.GetAppointmentsForDoctorResponse{Appointments: aptProtos}, nil
}

// mapUseCaseErrorToGrpcStatus converts use case errors to gRPC status errors.
func mapUseCaseErrorToGrpcStatus(err error) error {
	var ucErr *core_usecase.UseCaseError
	if errors.As(err, &ucErr) {
		switch ucErr.Type {
		case core_usecase.ErrNotFound:
			return status.Error(codes.NotFound, ucErr.Message)
		case core_usecase.ErrInvalidInput:
			return status.Error(codes.InvalidArgument, ucErr.Message)
		case core_usecase.ErrConflict:
			return status.Error(codes.AlreadyExists, ucErr.Message) // Or codes.FailedPrecondition / codes.Aborted
		case core_usecase.ErrInternal:
			return status.Error(codes.Internal, ucErr.Message)
		default:
			// Log the unknown type for debugging
			// log.Printf("Unknown use case error type: %v, message: %s", ucErr.Type, ucErr.Message)
			return status.Error(codes.Unknown, fmt.Sprintf("unknown use case error type: %v", ucErr.Type))
		}
	}
	// If it's not a UseCaseError, treat it as an internal server error.
	// Log the original error for debugging.
	// log.Printf("Unhandled internal error: %v", err) // Replace with proper logging
	return status.Error(codes.Internal, "an unexpected error occurred")
}
