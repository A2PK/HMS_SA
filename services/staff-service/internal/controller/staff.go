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

	coreUseCase "golang-microservices-boilerplate/pkg/core/usecase"
	pb "golang-microservices-boilerplate/proto/staff-service"
	"golang-microservices-boilerplate/services/staff-service/internal/usecase"
)

// StaffServer defines the interface for the gRPC service handler.
type StaffServer interface {
	pb.StaffServiceServer // Embed generated interface
}

// Ensure staffServer implements StaffServer.
var _ StaffServer = (*staffServer)(nil)

// --- gRPC Server Implementation ---

type staffServer struct {
	pb.UnimplementedStaffServiceServer
	uc     usecase.StaffUseCase
	mapper Mapper // Use interface
}

// NewStaffServer creates a new gRPC server instance.
func NewStaffServer(uc usecase.StaffUseCase, mapper Mapper) StaffServer {
	return &staffServer{
		uc:     uc,
		mapper: mapper,
	}
}

// RegisterStaffServiceServer registers the staff service implementation with the gRPC server.
func RegisterStaffServiceServer(s *grpc.Server, uc usecase.StaffUseCase, mapper Mapper) {
	server := NewStaffServer(uc, mapper)
	pb.RegisterStaffServiceServer(s, server)
}

// --- Staff Management RPCs ---

// AddStaff implements the corresponding gRPC method.
func (s *staffServer) AddStaff(ctx context.Context, req *pb.AddStaffRequest) (*pb.AddStaffResponse, error) {
	var dob *time.Time
	if req.DateOfBirth != nil && req.DateOfBirth.IsValid() {
		dobValue := req.DateOfBirth.AsTime()
		dob = &dobValue
	}

	staffEntity, err := s.uc.AddStaff(ctx, req.FirstName, req.LastName, dob, req.PhoneNumber, req.Address, req.RoleId, req.StatusId, req.Specialization, req.NurseType)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	staffProto, err := s.mapper.EntityToProto(staffEntity)
	if err != nil {
		// Log internal mapping error
		fmt.Printf("Error mapping staff entity to proto: %v\n", err) // Replace with proper logging
		return nil, status.Error(codes.Internal, "failed to process staff data")
	}

	return &pb.AddStaffResponse{Staff: staffProto}, nil
}

// GetStaffDetails implements the corresponding gRPC method.
func (s *staffServer) GetStaffDetails(ctx context.Context, req *pb.GetStaffDetailsRequest) (*pb.GetStaffDetailsResponse, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	staffEntity, err := s.uc.GetStaffDetails(ctx, staffID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	staffProto, err := s.mapper.EntityToProto(staffEntity)
	if err != nil {
		fmt.Printf("Error mapping staff entity to proto: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to process staff data")
	}

	return &pb.GetStaffDetailsResponse{Staff: staffProto}, nil
}

// UpdateStaffDetails implements the corresponding gRPC method.
func (s *staffServer) UpdateStaffDetails(ctx context.Context, req *pb.UpdateStaffDetailsRequest) (*pb.UpdateStaffDetailsResponse, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	var dob *time.Time
	if req.DateOfBirth != nil && req.DateOfBirth.IsValid() {
		dobValue := req.DateOfBirth.AsTime()
		dob = &dobValue
	}

	staffEntity, err := s.uc.UpdateStaffDetails(ctx, staffID, req.FirstName, req.LastName, dob, req.PhoneNumber, req.Address, req.Specialization, req.NurseType)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	staffProto, err := s.mapper.EntityToProto(staffEntity)
	if err != nil {
		fmt.Printf("Error mapping staff entity to proto: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to process staff data")
	}

	return &pb.UpdateStaffDetailsResponse{Staff: staffProto}, nil
}

// --- Restored RPC Implementations (Needing Review) ---

// UpdateStaffSchedule implements the corresponding gRPC method.
func (s *staffServer) UpdateStaffSchedule(ctx context.Context, req *pb.UpdateStaffScheduleRequest) (*emptypb.Empty, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	// Pass the TaskProto slice directly to the use case
	err = s.uc.UpdateStaffSchedule(ctx, staffID, req.TasksToSchedule)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

// SetStaffAvailability implements the corresponding gRPC method (now maps to SetStaffStatus).
func (s *staffServer) SetStaffAvailability(ctx context.Context, req *pb.SetStaffAvailabilityRequest) (*emptypb.Empty, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	// Call the renamed use case method
	err = s.uc.SetStaffStatus(ctx, staffID, req.StatusId)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

// GetDoctorAvailability implements the corresponding gRPC method.
func (s *staffServer) GetDoctorAvailability(ctx context.Context, req *pb.GetDoctorAvailabilityRequest) (*pb.GetDoctorAvailabilityResponse, error) {
	var doctorIDPtr *uuid.UUID
	if req.DoctorId != "" {
		parsedID, err := uuid.Parse(req.DoctorId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid doctor ID format: %v", err)
		}
		doctorIDPtr = &parsedID
	}

	startTime := time.Time{}
	if req.StartTime != nil && req.StartTime.IsValid() {
		startTime = req.StartTime.AsTime()
	}
	endTime := time.Time{}
	if req.EndTime != nil && req.EndTime.IsValid() {
		endTime = req.EndTime.AsTime()
	}

	if startTime.IsZero() || endTime.IsZero() || endTime.Before(startTime) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time range provided")
	}

	// Use case now returns []*pb.GetDoctorAvailabilityResponse_TimeSlot directly
	availableSlotsProto, err := s.uc.GetDoctorAvailability(ctx, doctorIDPtr, startTime, endTime)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	// No further mapping needed as use case returns the correct proto type
	return &pb.GetDoctorAvailabilityResponse{AvailableSlots: availableSlotsProto}, nil
}

// AssignTask implements the corresponding gRPC method.
func (s *staffServer) AssignTask(ctx context.Context, req *pb.AssignTaskRequest) (*emptypb.Empty, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	var startTime, endTime *time.Time
	if req.StartTime != nil && req.StartTime.IsValid() {
		st := req.StartTime.AsTime()
		startTime = &st
	}
	if req.EndTime != nil && req.EndTime.IsValid() {
		et := req.EndTime.AsTime()
		endTime = &et
	}

	_, err = s.uc.AssignTask(ctx, staffID, req.Title, req.Description, req.Priority, startTime, endTime, req.StatusId)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	// Returning Empty as per proto definition, could return created task if needed
	return &emptypb.Empty{}, nil
}

// TrackWorkload implements the corresponding gRPC method.
func (s *staffServer) TrackWorkload(ctx context.Context, req *pb.TrackWorkloadRequest) (*pb.TrackWorkloadResponse, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staff ID format: %v", err)
	}

	workloadEntities, err := s.uc.TrackWorkload(ctx, staffID)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}

	workloadProto, err := s.mapper.TasksToProto(workloadEntities)
	if err != nil {
		fmt.Printf("Error mapping workload tasks to proto: %v\n", err)
		return nil, status.Error(codes.Internal, "failed to process workload data")
	}

	return &pb.TrackWorkloadResponse{Workload: workloadProto}, nil
}

// --- Lookup Table RPCs ---

// AddStaffRole implements the gRPC method.
func (s *staffServer) AddStaffRole(ctx context.Context, req *pb.AddStaffRoleRequest) (*pb.AddStaffRoleResponse, error) {
	roleEntity, err := s.uc.AddStaffRole(ctx, req.Name, req.Description)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.AddStaffRoleResponse{Role: s.mapper.RoleToProto(roleEntity)}, nil
}

// ListStaffRoles implements the gRPC method.
func (s *staffServer) ListStaffRoles(ctx context.Context, req *pb.ListStaffRolesRequest) (*pb.ListStaffRolesResponse, error) {
	roleEntities, err := s.uc.ListStaffRoles(ctx)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.ListStaffRolesResponse{Roles: s.mapper.RolesToProto(roleEntities)}, nil
}

// AddStaffStatus implements the gRPC method.
func (s *staffServer) AddStaffStatus(ctx context.Context, req *pb.AddStaffStatusRequest) (*pb.AddStaffStatusResponse, error) {
	statusEntity, err := s.uc.AddStaffStatus(ctx, req.Name, req.Description)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.AddStaffStatusResponse{Status: s.mapper.StatusToProto(statusEntity)}, nil
}

// ListStaffStatuses implements the gRPC method.
func (s *staffServer) ListStaffStatuses(ctx context.Context, req *pb.ListStaffStatusesRequest) (*pb.ListStaffStatusesResponse, error) {
	statusEntities, err := s.uc.ListStaffStatuses(ctx)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.ListStaffStatusesResponse{Statuses: s.mapper.StatusesToProto(statusEntities)}, nil
}

// AddTaskStatus implements the gRPC method.
func (s *staffServer) AddTaskStatus(ctx context.Context, req *pb.AddTaskStatusRequest) (*pb.AddTaskStatusResponse, error) {
	statusEntity, err := s.uc.AddTaskStatus(ctx, req.Name, req.Description)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.AddTaskStatusResponse{Status: s.mapper.TaskStatusToProto(statusEntity)}, nil
}

// ListTaskStatuses implements the gRPC method.
func (s *staffServer) ListTaskStatuses(ctx context.Context, req *pb.ListTaskStatusesRequest) (*pb.ListTaskStatusesResponse, error) {
	statusEntities, err := s.uc.ListTaskStatuses(ctx)
	if err != nil {
		return nil, mapUseCaseErrorToGrpcStatus(err)
	}
	return &pb.ListTaskStatusesResponse{Statuses: s.mapper.TaskStatusesToProto(statusEntities)}, nil
}

// --- Helper Functions ---

// mapUseCaseErrorToGrpcStatus converts use case errors to gRPC status errors.
func mapUseCaseErrorToGrpcStatus(err error) error {
	var ucErr *coreUseCase.UseCaseError // Use aliased import
	if errors.As(err, &ucErr) {
		switch ucErr.Type {
		case coreUseCase.ErrNotFound:
			return status.Error(codes.NotFound, ucErr.Message)
		case coreUseCase.ErrInvalidInput:
			return status.Error(codes.InvalidArgument, ucErr.Message)
		case coreUseCase.ErrConflict:
			return status.Error(codes.AlreadyExists, ucErr.Message)
		case coreUseCase.ErrInternal:
			return status.Error(codes.Internal, ucErr.Message)
		default:
			return status.Error(codes.Unknown, fmt.Sprintf("unknown use case error type: %v", ucErr.Type))
		}
	}
	// Log unexpected non-UseCaseError errors
	fmt.Printf("Unexpected error type encountered: %T - %v\n", err, err) // Replace with proper logging
	return status.Error(codes.Internal, "an unexpected internal error occurred")
}
