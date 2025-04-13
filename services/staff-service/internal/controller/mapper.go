package controller

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "golang-microservices-boilerplate/proto/staff-service"
	"golang-microservices-boilerplate/services/staff-service/internal/entity"
)

// Mapper defines the interface for mapping between gRPC proto messages and internal staff types.
// TODO: Consider code generation tools (like protoc-gen-go-copy) for complex mappings.
type Mapper interface {
	// Staff Mappings
	EntityToProto(staff *entity.Staff) (*pb.Staff, error)

	// Schedule & Task Mappings
	ScheduleEntryToProto(entry *entity.ScheduleEntry) (*pb.ScheduleEntryProto, error)
	ScheduleEntriesToProto(entries []entity.ScheduleEntry) ([]*pb.ScheduleEntryProto, error)
	TaskToProto(task *entity.Task) (*pb.TaskProto, error)
	TasksToProto(tasks []*entity.Task) ([]*pb.TaskProto, error)

	// Lookup Mappings (Entity -> Proto)
	RoleToProto(role *entity.StaffRole) *pb.StaffRoleProto
	RolesToProto(roles []*entity.StaffRole) []*pb.StaffRoleProto
	StatusToProto(status *entity.StaffStatus) *pb.StaffStatusProto
	StatusesToProto(statuses []*entity.StaffStatus) []*pb.StaffStatusProto
	TaskStatusToProto(status *entity.TaskStatus) *pb.TaskStatusProto
	TaskStatusesToProto(statuses []*entity.TaskStatus) []*pb.TaskStatusProto

	// Proto -> Entity (Only where needed for requests, often handled directly in use case)
	// Example: ProtoAssignRequestToEntity might not be needed if use case takes direct args.
	// ProtoScheduleEntriesToEntity removed (was complex, likely handled differently now)
}

// Ensure StaffMapper implements Mapper interface.
var _ Mapper = (*StaffMapper)(nil)

// StaffMapper handles mapping between gRPC proto messages and internal staff types.
type StaffMapper struct{}

// NewStaffMapper creates a new instance of StaffMapper.
func NewStaffMapper() Mapper { // Return interface type
	return &StaffMapper{}
}

// --- Staff Mappings ---

// EntityToProto converts an entity.Staff to a proto.Staff.
func (m *StaffMapper) EntityToProto(staff *entity.Staff) (*pb.Staff, error) {
	if staff == nil {
		return nil, errors.New("cannot map nil staff entity to proto")
	}

	scheduleProto, err := m.ScheduleEntriesToProto(staff.Schedule) // Maps []entity.ScheduleEntry
	if err != nil {
		return nil, fmt.Errorf("failed to map schedule: %w", err)
	}

	var dobProto *timestamppb.Timestamp
	if !staff.DateOfBirth.IsZero() {
		dobProto = timestamppb.New(staff.DateOfBirth)
	}

	return &pb.Staff{
		Id:             staff.ID.String(),
		FirstName:      staff.FirstName,
		LastName:       staff.LastName,
		DateOfBirth:    dobProto,
		PhoneNumber:    staff.PhoneNumber,
		Address:        staff.Address,
		RoleId:         staff.RoleID,   // String FK
		StatusId:       staff.StatusID, // String FK
		Specialization: staff.Specialization,
		NurseType:      staff.NurseType,
		Schedule:       scheduleProto, // Updated field type
		CreatedAt:      timestamppb.New(staff.CreatedAt),
		UpdatedAt:      timestamppb.New(staff.UpdatedAt),
	}, nil
}

// --- Schedule & Task Mappings ---

// ScheduleEntryToProto converts entity.ScheduleEntry to proto.ScheduleEntryProto.
func (m *StaffMapper) ScheduleEntryToProto(entry *entity.ScheduleEntry) (*pb.ScheduleEntryProto, error) {
	if entry == nil {
		return nil, errors.New("cannot map nil schedule entry to proto")
	}
	taskProto, err := m.TaskToProto(&entry.Task) // Map nested Task
	if err != nil {
		// Decide how to handle error: return error, return partial, or skip?
		// Returning error for now, as ScheduleEntry must have a Task.
		return nil, fmt.Errorf("failed to map task for schedule entry: %w", err)
	}

	return &pb.ScheduleEntryProto{
		StaffId: entry.StaffID.String(),
		Task:    taskProto,
	}, nil
}

// ScheduleEntriesToProto converts a slice of entity.ScheduleEntry to proto.
func (m *StaffMapper) ScheduleEntriesToProto(entries []entity.ScheduleEntry) ([]*pb.ScheduleEntryProto, error) {
	protos := make([]*pb.ScheduleEntryProto, 0, len(entries))
	for i := range entries { // Iterate by index to get pointer
		proto, err := m.ScheduleEntryToProto(&entries[i])
		if err != nil {
			// Log or handle error - skipping entry for now
			fmt.Printf("Error mapping schedule entry %v: %v\n", entries[i], err)
			continue
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// TaskToProto converts entity.Task to proto.TaskProto.
func (m *StaffMapper) TaskToProto(task *entity.Task) (*pb.TaskProto, error) {
	if task == nil {
		return nil, errors.New("cannot map nil task to proto")
	}
	// Note: Assumes task.Status (TaskStatus entity) is preloaded if needed for display
	// Here we only map the StatusID string

	return &pb.TaskProto{
		Id:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Priority:    int32(task.Priority),
		StartTime:   timestamppb.New(task.StartTime),
		EndTime:     timestamppb.New(task.EndTime),
		StatusId:    task.StatusID, // String FK
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil
}

// TasksToProto converts a slice of *entity.Task to proto.
func (m *StaffMapper) TasksToProto(tasks []*entity.Task) ([]*pb.TaskProto, error) {
	protos := make([]*pb.TaskProto, 0, len(tasks))
	for _, task := range tasks {
		proto, err := m.TaskToProto(task)
		if err != nil {
			// Log or handle error - skipping task for now
			fmt.Printf("Error mapping task %v: %v\n", task.ID, err)
			continue
		}
		protos = append(protos, proto)
	}
	return protos, nil
}

// --- Lookup Table Mappings ---

// RoleToProto maps entity.StaffRole to proto.
func (m *StaffMapper) RoleToProto(role *entity.StaffRole) *pb.StaffRoleProto {
	if role == nil {
		return nil
	}
	return &pb.StaffRoleProto{
		Name:        role.Name,
		Description: role.Description,
	}
}

// RolesToProto maps a slice of entity.StaffRole to proto.
func (m *StaffMapper) RolesToProto(roles []*entity.StaffRole) []*pb.StaffRoleProto {
	protos := make([]*pb.StaffRoleProto, len(roles))
	for i, role := range roles {
		protos[i] = m.RoleToProto(role)
	}
	return protos
}

// StatusToProto maps entity.StaffStatus to proto.
func (m *StaffMapper) StatusToProto(status *entity.StaffStatus) *pb.StaffStatusProto {
	if status == nil {
		return nil
	}
	return &pb.StaffStatusProto{
		Name:        status.Name,
		Description: status.Description,
	}
}

// StatusesToProto maps a slice of entity.StaffStatus to proto.
func (m *StaffMapper) StatusesToProto(statuses []*entity.StaffStatus) []*pb.StaffStatusProto {
	protos := make([]*pb.StaffStatusProto, len(statuses))
	for i, status := range statuses {
		protos[i] = m.StatusToProto(status)
	}
	return protos
}

// TaskStatusToProto maps entity.TaskStatus to proto.
func (m *StaffMapper) TaskStatusToProto(status *entity.TaskStatus) *pb.TaskStatusProto {
	if status == nil {
		return nil
	}
	return &pb.TaskStatusProto{
		Name:        status.Name,
		Description: status.Description,
	}
}

// TaskStatusesToProto maps a slice of entity.TaskStatus to proto.
func (m *StaffMapper) TaskStatusesToProto(statuses []*entity.TaskStatus) []*pb.TaskStatusProto {
	protos := make([]*pb.TaskStatusProto, len(statuses))
	for i, status := range statuses {
		protos[i] = m.TaskStatusToProto(status)
	}
	return protos
}

// Removed Proto -> Entity mappings for StaffType, ScheduleEntry, AssignTaskRequest
// These are typically handled by extracting necessary fields in the use case or handler.
