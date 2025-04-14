package repository

import (
	"context"
	"time"

	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/staff-service/internal/entity"

	"github.com/google/uuid"
)

// StaffRepository defines the interface for staff data persistence.
// It embeds the BaseRepository for standard CRUD on entity.Staff.
// Implementations should handle preloading of related entities (Role, Status, Schedule.Task) where appropriate.
type StaffRepository interface {
	coreRepository.BaseRepository[entity.Staff]

	// FindAvailableDoctors retrieves staff (typically doctors by role) who are considered available
	// within a given time range. Implementation needs to consider Staff.StatusID and existing Task start/end times.
	FindAvailableDoctors(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entity.Staff, error)

	// AddScheduleEntries creates new Task entities and links them to the specified staff member via ScheduleEntry records.
	AddScheduleEntries(ctx context.Context, staffID uuid.UUID, tasks []*entity.Task) error

	// UpdateStatus updates the StatusID field for a specific staff member.
	UpdateStatus(ctx context.Context, staffID uuid.UUID, statusID string) error

	// AssignTaskToStaff creates a Task and links it to the staff member via a ScheduleEntry.
	AssignTaskToStaff(ctx context.Context, staffID uuid.UUID, task *entity.Task) error

	// FindTasksByStaffID retrieves all tasks associated with a staff member via their schedule.
	FindTasksByStaffID(ctx context.Context, staffID uuid.UUID) ([]*entity.Task, error)

	// FindByID (from embedded BaseRepository) should handle preloading
	// FindByID(ctx context.Context, id uuid.UUID) (*entity.Staff, error)
}

// Add TaskRepository interface definition
type TaskRepository interface {
	coreRepository.BaseRepository[entity.Task]
	// Add any Task-specific methods here if needed beyond BaseRepository
}

// --- Interfaces for Lookup Tables --- //

// StaffRoleRepository defines the interface for staff role data persistence.
type StaffRoleRepository interface {
	Create(ctx context.Context, role *entity.StaffRole) error
	FindByName(ctx context.Context, name string) (*entity.StaffRole, error)
	// ListAll returns a simple slice, as pagination based on BaseEntity fields isn't applicable.
	ListAll(ctx context.Context) ([]*entity.StaffRole, error)
	// Update and Delete might be needed depending on requirements.
}

// StaffStatusRepository defines the interface for staff status data persistence.
type StaffStatusRepository interface {
	Create(ctx context.Context, status *entity.StaffStatus) error
	FindByName(ctx context.Context, name string) (*entity.StaffStatus, error)
	// ListAll returns a simple slice.
	ListAll(ctx context.Context) ([]*entity.StaffStatus, error)
}

// TaskStatusRepository defines the interface for task status data persistence.
type TaskStatusRepository interface {
	Create(ctx context.Context, status *entity.TaskStatus) error
	FindByName(ctx context.Context, name string) (*entity.TaskStatus, error)
	// ListAll returns a simple slice.
	ListAll(ctx context.Context) ([]*entity.TaskStatus, error)
}

// Note: Consider custom error types like ErrStaffNotFound, ErrRoleNotFound, etc.
