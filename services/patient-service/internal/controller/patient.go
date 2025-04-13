package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	core_usecase "golang-microservices-boilerplate/pkg/core/usecase"
	pb "golang-microservices-boilerplate/proto/patient-service"
	"golang-microservices-boilerplate/services/patient-service/internal/usecase"
)

// PatientServer defines the interface for the gRPC service handler.
// This corresponds to the pb.PatientServiceServer interface but allows for dependency injection.
type PatientServer interface {
	pb.PatientServiceServer // Embed the generated interface
	// Add any other methods specific to the server lifecycle if needed
}

// Ensure patientServer implements PatientServer interface (and pb.PatientServiceServer).
var _ PatientServer = (*patientServer)(nil)

// --- gRPC Server Implementation ---

type patientServer struct {
	pb.UnimplementedPatientServiceServer
	uc     usecase.PatientUseCase
	mapper Mapper // Use the Mapper interface
}

// NewPatientServer creates a new gRPC server instance.
// Accepts Mapper interface and returns PatientServer interface.
func NewPatientServer(uc usecase.PatientUseCase, mapper Mapper) PatientServer {
	return &patientServer{
		uc:     uc,
		mapper: mapper, // Inject mapper
	}
}

// RegisterPatientServiceServer registers the patient service implementation with the gRPC server.
// Accepts use case and mapper to create the server.
func RegisterPatientServiceServer(s *grpc.Server, uc usecase.PatientUseCase, mapper Mapper) {
	server := NewPatientServer(uc, mapper) // Pass mapper
	pb.RegisterPatientServiceServer(s, server)
}

// RegisterPatient implements the corresponding gRPC method.
func (s *patientServer) RegisterPatient(ctx context.Context, req *pb.RegisterPatientRequest) (*pb.RegisterPatientResponse, error) {
	// Use case already accepts the proto request.
	patientEntity, err := s.uc.RegisterPatient(ctx, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	patientProto, err := s.mapper.EntityToProto(patientEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.RegisterPatientResponse{Patient: patientProto}, nil
}

// GetPatientDetails implements the corresponding gRPC method.
func (s *patientServer) GetPatientDetails(ctx context.Context, req *pb.GetPatientDetailsRequest) (*pb.GetPatientDetailsResponse, error) {
	patientID, err := uuid.Parse(req.PatientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid patient ID format: %v", err)
	}

	patientEntity, err := s.uc.GetPatientDetails(ctx, patientID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	patientProto, err := s.mapper.EntityToProto(patientEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.GetPatientDetailsResponse{Patient: patientProto}, nil
}

// UpdatePatientDetails implements the corresponding gRPC method.
func (s *patientServer) UpdatePatientDetails(ctx context.Context, req *pb.UpdatePatientDetailsRequest) (*pb.UpdatePatientDetailsResponse, error) {
	patientID, err := uuid.Parse(req.PatientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid patient ID format: %v", err)
	}

	// Use case already accepts the proto request.
	patientEntity, err := s.uc.UpdatePatientDetails(ctx, patientID, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	patientProto, err := s.mapper.EntityToProto(patientEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map result to proto: %v", err)
	}

	return &pb.UpdatePatientDetailsResponse{Patient: patientProto}, nil
}

// AddMedicalRecord implements the corresponding gRPC method.
func (s *patientServer) AddMedicalRecord(ctx context.Context, req *pb.AddMedicalRecordRequest) (*emptypb.Empty, error) {
	patientID, err := uuid.Parse(req.PatientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid patient ID format: %v", err)
	}

	// Use case already accepts the proto request.
	err = s.uc.AddMedicalRecord(ctx, patientID, req)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

// GetPatientMedicalHistory implements the corresponding gRPC method.
func (s *patientServer) GetPatientMedicalHistory(ctx context.Context, req *pb.GetPatientMedicalHistoryRequest) (*pb.GetPatientMedicalHistoryResponse, error) {
	patientID, err := uuid.Parse(req.PatientId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid patient ID format: %v", err)
	}

	medicalHistoryEntities, err := s.uc.GetPatientMedicalHistory(ctx, patientID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	medicalHistoryProto, err := s.mapper.MedicalRecordsToProto(medicalHistoryEntities)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map results to proto: %v", err)
	}

	return &pb.GetPatientMedicalHistoryResponse{MedicalHistory: medicalHistoryProto}, nil
}

// mapUseCaseErrorToGrpcStatus converts use case errors to gRPC status errors.
// (This is the same helper function as in appointment controller)
func mapUseCaseErrorToGrpcStatus(err error) error {
	var ucErr *core_usecase.UseCaseError
	if errors.As(err, &ucErr) {
		switch ucErr.Type {
		case core_usecase.ErrNotFound:
			return status.Error(codes.NotFound, ucErr.Message)
		case core_usecase.ErrInvalidInput:
			return status.Error(codes.InvalidArgument, ucErr.Message)
		case core_usecase.ErrConflict:
			return status.Error(codes.AlreadyExists, ucErr.Message) // Or FailedPrecondition
		case core_usecase.ErrInternal:
			return status.Error(codes.Internal, ucErr.Message)
		default:
			return status.Error(codes.Unknown, fmt.Sprintf("unknown use case error type: %v", ucErr.Type))
		}
	}
	return status.Error(codes.Internal, "an unexpected error occurred")
}
