package usecase

import (
	"context"
	"time"

	pb "golang-microservices-boilerplate/proto/staff-service"
	"golang-microservices-boilerplate/services/staff-service/internal/entity"

	"github.com/google/uuid"
)

// StaffUseCase defines the interface for staff and related lookup management operations.
type StaffUseCase interface {
	// === Staff Management ===

	// AddStaff adds a new staff member.
	// Takes parameters aligned with pb.AddStaffRequest.
	AddStaff(ctx context.Context, firstName, lastName string, dob *time.Time, phone, address, roleID, statusID, specialization, nurseType string) (*entity.Staff, error)

	// GetStaffDetails retrieves details of a specific staff member by ID, including schedule and tasks.
	GetStaffDetails(ctx context.Context, staffID uuid.UUID) (*entity.Staff, error)

	// UpdateStaffDetails updates information for an existing staff member.
	// Takes parameters aligned with pb.UpdateStaffDetailsRequest.
	UpdateStaffDetails(ctx context.Context, staffID uuid.UUID, firstName, lastName string, dob *time.Time, phone, address, specialization, nurseType string) (*entity.Staff, error)

	// UpdateStaffSchedule updates the schedule for a staff member by creating new tasks and schedule entries.
	// Input tasks are expected to be DTOs or similar, not raw entities.
	UpdateStaffSchedule(ctx context.Context, staffID uuid.UUID, tasksToSchedule []*pb.TaskProto) error // TODO: Define a suitable DTO instead of pb.TaskProto?

	// SetStaffStatus updates the status (StatusID) of a staff member.
	SetStaffStatus(ctx context.Context, staffID uuid.UUID, statusID string) error

	// GetDoctorAvailability calculates and retrieves available time slots for doctors.
	// Input doctorID is optional (nil checks all doctors).
	// Return type needs finalization (e.g., list of TimeSlot DTOs).
	GetDoctorAvailability(ctx context.Context, doctorID *uuid.UUID, startTime time.Time, endTime time.Time) ([]*pb.GetDoctorAvailabilityResponse_TimeSlot, error)

	// AssignTask creates a task and assigns it to a staff member via a schedule entry.
	// Input parameters align with pb.AssignTaskRequest.
	AssignTask(ctx context.Context, staffID uuid.UUID, title, description string, priority int32, startTime, endTime *time.Time, statusID string) (*entity.Task, error)

	// TrackWorkload retrieves the current workload (tasks) for a staff member.
	TrackWorkload(ctx context.Context, staffID uuid.UUID) ([]*entity.Task, error)

	// ListStaff retrieves a list of staff members, optionally filtering by role and/or status.
	ListStaff(ctx context.Context, req *pb.ListStaffRequest) ([]*entity.Staff, error)

	// ListTasks retrieves a list of all tasks, optionally filtered, ordered by creation time descending.
	ListTasks(ctx context.Context, req *pb.ListTasksRequest) ([]*entity.Task, error)

	// === Lookup Table Management ===

	// Staff Roles
	AddStaffRole(ctx context.Context, name, description string) (*entity.StaffRole, error)
	ListStaffRoles(ctx context.Context) ([]*entity.StaffRole, error)

	// Staff Statuses
	AddStaffStatus(ctx context.Context, name, description string) (*entity.StaffStatus, error)
	ListStaffStatuses(ctx context.Context) ([]*entity.StaffStatus, error)

	// Task Statuses
	AddTaskStatus(ctx context.Context, name, description string) (*entity.TaskStatus, error)
	ListTaskStatuses(ctx context.Context) ([]*entity.TaskStatus, error)
}
