package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "golang-microservices-boilerplate/proto/staff-service"

	coreLogger "golang-microservices-boilerplate/pkg/core/logger"
	coreUseCase "golang-microservices-boilerplate/pkg/core/usecase"

	// coreDTO "golang-microservices-boilerplate/pkg/core/dto"

	"golang-microservices-boilerplate/services/staff-service/internal/entity"
	"golang-microservices-boilerplate/services/staff-service/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Ensure implementation satisfies the interface
var _ StaffUseCase = (*staffUseCaseImpl)(nil)

// staffUseCaseImpl implements the StaffUseCase interface.
type staffUseCaseImpl struct {
	// Remove BaseUseCaseImpl embedding as we handle methods specifically
	staffRepo       repository.StaffRepository // Specific repository interface for Staff
	staffRoleRepo   repository.StaffRoleRepository
	staffStatusRepo repository.StaffStatusRepository
	taskStatusRepo  repository.TaskStatusRepository
	logger          coreLogger.Logger
}

// NewStaffUseCase creates a new instance of staffUseCaseImpl.
func NewStaffUseCase(
	staffRepo repository.StaffRepository,
	staffRoleRepo repository.StaffRoleRepository,
	staffStatusRepo repository.StaffStatusRepository,
	taskStatusRepo repository.TaskStatusRepository,
	logger coreLogger.Logger,
) StaffUseCase {
	return &staffUseCaseImpl{
		staffRepo:       staffRepo,
		staffRoleRepo:   staffRoleRepo,
		staffStatusRepo: staffStatusRepo,
		taskStatusRepo:  taskStatusRepo,
		logger:          logger,
	}
}

// --- Staff Management Implementations ---

// AddStaff handles creating a new staff member.
func (uc *staffUseCaseImpl) AddStaff(ctx context.Context, firstName, lastName string, dob *time.Time, phone, address, roleID, statusID, specialization, nurseType string) (*entity.Staff, error) {
	uc.logger.Info("Adding new staff", "firstName", firstName, "lastName", lastName, "roleID", roleID, "statusID", statusID)

	// Validation
	if firstName == "" || lastName == "" || phone == "" || roleID == "" || statusID == "" || dob == nil || dob.IsZero() {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "missing required staff information")
	}

	// Validate existence of roleID and statusID
	_, err := uc.staffRoleRepo.FindByName(ctx, roleID)
	if err != nil {
		uc.logger.Warn("Invalid RoleID provided", "roleID", roleID, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid role ID: %s", roleID))
	}
	_, err = uc.staffStatusRepo.FindByName(ctx, statusID)
	if err != nil {
		uc.logger.Warn("Invalid StatusID provided", "statusID", statusID, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid status ID: %s", statusID))
	}

	// Create entity using domain constructor
	staff := entity.NewStaff(firstName, lastName, phone, address, *dob, roleID, statusID, specialization, nurseType)

	// Use the specific repository's Create method
	err = uc.staffRepo.Create(ctx, staff) // GormBaseRepository provides Create
	if err != nil {
		uc.logger.Error("Failed to save staff", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to add staff member")
	}

	uc.logger.Info("Staff added successfully", "staffID", staff.ID.String())
	// Refetch to preload relations for the response
	return uc.staffRepo.FindByID(ctx, staff.ID)
}

// GetStaffDetails retrieves staff details by ID.
func (uc *staffUseCaseImpl) GetStaffDetails(ctx context.Context, staffID uuid.UUID) (*entity.Staff, error) {
	uc.logger.Info("Getting staff details", "staffID", staffID.String())
	if staffID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID")
	}

	// Use the specific repository's FindByID which handles preloading
	staff, err := uc.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, errors.New("staff not found")) {
			uc.logger.Warn("Staff not found", "staffID", staffID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "staff not found")
		}
		uc.logger.Error("Failed to get staff details", "staffID", staffID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to get staff details")
	}
	return staff, nil
}

// UpdateStaffDetails updates existing staff details.
func (uc *staffUseCaseImpl) UpdateStaffDetails(ctx context.Context, staffID uuid.UUID, firstName, lastName string, dob *time.Time, phone, address, specialization, nurseType string) (*entity.Staff, error) {
	uc.logger.Info("Updating staff details", "staffID", staffID.String())
	if staffID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID")
	}

	// Fetch existing staff - Use the base FindByID without preloads for the update target
	staff, err := uc.staffRepo.FindByID(ctx, staffID) // staffRepo embeds the base repo, call FindByID directly
	if err != nil {
		if errors.Is(err, errors.New("staff not found")) {
			uc.logger.Warn("Staff not found for update", "staffID", staffID.String())
			return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "staff not found")
		}
		uc.logger.Error("Failed to find staff for update", "staffID", staffID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to find staff for update")
	}

	// Apply updates using entity method
	var dobValue time.Time
	if dob != nil {
		dobValue = *dob
	}
	staff.UpdateDetails(firstName, lastName, phone, address, dobValue, specialization, nurseType)

	// Use repository's Update method
	err = uc.staffRepo.Update(ctx, staff)
	if err != nil {
		uc.logger.Error("Failed to update staff", "staffID", staffID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to update staff details")
	}

	uc.logger.Info("Staff details updated successfully", "staffID", staffID.String())
	// Refetch using the specific FindByID to get preloaded relations for the response
	return uc.staffRepo.FindByID(ctx, staffID)
}

// UpdateStaffSchedule creates new tasks and links them.
func (uc *staffUseCaseImpl) UpdateStaffSchedule(ctx context.Context, staffID uuid.UUID, tasksToSchedule []*pb.TaskProto) error {
	uc.logger.Info("Updating staff schedule by adding tasks", "staffID", staffID.String(), "taskCount", len(tasksToSchedule))
	if staffID == uuid.Nil {
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID")
	}

	// Map []*pb.TaskProto to []*entity.Task and validate
	tasks := make([]*entity.Task, 0, len(tasksToSchedule))
	for _, pt := range tasksToSchedule {
		if pt.StartTime == nil || pt.EndTime == nil || pt.StartTime.AsTime().After(pt.EndTime.AsTime()) || pt.StatusId == "" {
			uc.logger.Warn("Invalid task data in schedule update", "staffID", staffID.String(), "taskTitle", pt.Title)
			return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid data for task '%s' (time range or status ID)", pt.Title))
		}
		// Validate existence of StatusId
		_, err := uc.taskStatusRepo.FindByName(ctx, pt.StatusId)
		if err != nil {
			uc.logger.Warn("Invalid TaskStatusID provided for task", "taskTitle", pt.Title, "statusID", pt.StatusId, "error", err)
			return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid status ID '%s' for task '%s'", pt.StatusId, pt.Title))
		}

		tasks = append(tasks, &entity.Task{
			Title:       pt.Title,
			Description: pt.Description,
			Priority:    int(pt.Priority),
			StartTime:   pt.StartTime.AsTime(),
			EndTime:     pt.EndTime.AsTime(),
			StatusID:    pt.StatusId,
		})
	}

	err := uc.staffRepo.AddScheduleEntries(ctx, staffID, tasks)
	if err != nil {
		uc.logger.Error("Failed to add schedule entries in repo", "staffID", staffID.String(), "error", err)
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to update staff schedule")
	}

	uc.logger.Info("Staff schedule updated successfully by adding tasks", "staffID", staffID.String())
	return nil
}

// SetStaffStatus updates the staff's status.
func (uc *staffUseCaseImpl) SetStaffStatus(ctx context.Context, staffID uuid.UUID, statusID string) error {
	uc.logger.Info("Setting staff status", "staffID", staffID.String(), "statusID", statusID)
	if staffID == uuid.Nil || statusID == "" {
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID or status ID")
	}
	// Validate existence of statusID
	_, err := uc.staffStatusRepo.FindByName(ctx, statusID)
	if err != nil {
		uc.logger.Warn("Invalid StatusID provided for update", "statusID", statusID, "error", err)
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid status ID: %s", statusID))
	}

	err = uc.staffRepo.UpdateStatus(ctx, staffID, statusID)
	if err != nil {
		uc.logger.Error("Failed to set staff status in repo", "staffID", staffID.String(), "error", err)
		if errors.Is(err, errors.New("staff not found or status not changed")) {
			return coreUseCase.NewUseCaseError(coreUseCase.ErrNotFound, "staff not found")
		}
		return coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to set staff status")
	}
	uc.logger.Info("Staff status set successfully", "staffID", staffID.String())
	return nil
}

// GetDoctorAvailability retrieves doctor availability.
// TODO: Implement actual availability calculation based on task times.
func (uc *staffUseCaseImpl) GetDoctorAvailability(ctx context.Context, doctorID *uuid.UUID, startTime time.Time, endTime time.Time) ([]*pb.GetDoctorAvailabilityResponse_TimeSlot, error) {
	logFields := []interface{}{"startTime", startTime, "endTime", endTime}
	if doctorID != nil {
		logFields = append(logFields, "doctorID", doctorID.String())
	}
	uc.logger.Info("Getting doctor availability (placeholder implementation)", logFields...)

	if doctorID != nil && *doctorID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid doctor ID provided")
	}
	if startTime.IsZero() || endTime.IsZero() || endTime.Before(startTime) {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid time range")
	}

	// Using repo FindAvailableDoctors (which also needs refactoring for actual calculation)
	// For now, this gets active doctors but doesn't calculate free slots.
	doctors, err := uc.staffRepo.FindAvailableDoctors(ctx, startTime, endTime)
	if err != nil {
		uc.logger.Error("Failed to find active doctors", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to retrieve potential doctor availability")
	}

	// If a specific doctorID is requested, filter the list
	if doctorID != nil {
		filteredDoctors := []*entity.Staff{}
		for _, d := range doctors {
			if d.ID == *doctorID {
				filteredDoctors = append(filteredDoctors, d)
				break
			}
		}
		doctors = filteredDoctors // Replace original list with the filtered one (or empty)
	}

	// Placeholder: Return a single slot covering the entire requested range if any active doctor exists
	// The actual logic should calculate free slots based on each doctor's tasks.
	availableSlots := []*pb.GetDoctorAvailabilityResponse_TimeSlot{}
	if len(doctors) > 0 {
		uc.logger.Warn("GetDoctorAvailability returning placeholder availability - full range for first found doctor")
		availableSlots = append(availableSlots, &pb.GetDoctorAvailabilityResponse_TimeSlot{
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
		})
	}

	uc.logger.Info("Finished getting doctor availability (placeholder)", "slotCount", len(availableSlots))
	return availableSlots, nil
}

// AssignTask assigns a task.
func (uc *staffUseCaseImpl) AssignTask(ctx context.Context, staffID uuid.UUID, title, description string, priority int32, startTime, endTime *time.Time, statusID string) (*entity.Task, error) {
	uc.logger.Info("Assigning task", "staffID", staffID.String(), "title", title)
	if staffID == uuid.Nil || title == "" || statusID == "" || startTime == nil || endTime == nil || startTime.IsZero() || endTime.IsZero() || startTime.After(*endTime) {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid input for assigning task")
	}

	// Validate existence of staffID and statusID
	_, err := uc.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		uc.logger.Warn("Invalid StaffID provided for task assignment", "staffID", staffID, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid staff ID: %s", staffID))
	}
	_, err = uc.taskStatusRepo.FindByName(ctx, statusID)
	if err != nil {
		uc.logger.Warn("Invalid TaskStatusID provided for task assignment", "statusID", statusID, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, fmt.Sprintf("invalid task status ID: %s", statusID))
	}

	task := &entity.Task{
		Title:       title,
		Description: description,
		Priority:    int(priority),
		StartTime:   *startTime,
		EndTime:     *endTime,
		StatusID:    statusID,
	}

	err = uc.staffRepo.AssignTaskToStaff(ctx, staffID, task)
	if err != nil {
		uc.logger.Error("Failed to assign task in repo", "staffID", staffID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to assign task")
	}

	uc.logger.Info("Task assigned successfully", "staffID", staffID.String(), "taskID", task.ID.String())
	// Fetch status relation for the response
	taskStatus, findErr := uc.taskStatusRepo.FindByName(ctx, task.StatusID)
	if findErr == nil && taskStatus != nil {
		task.Status = *taskStatus // Dereference the pointer here
	} else if findErr != nil {
		// Log error if finding the status failed, but don't fail the whole operation
		uc.logger.Error("Failed to find task status after assignment", "statusID", task.StatusID, "error", findErr)
	}
	return task, nil
}

// TrackWorkload retrieves tasks for a staff member.
func (uc *staffUseCaseImpl) TrackWorkload(ctx context.Context, staffID uuid.UUID) ([]*entity.Task, error) {
	uc.logger.Info("Tracking workload", "staffID", staffID.String())
	if staffID == uuid.Nil {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "invalid staff ID")
	}

	tasks, err := uc.staffRepo.FindTasksByStaffID(ctx, staffID)
	if err != nil {
		uc.logger.Error("Failed to track workload in repo", "staffID", staffID.String(), "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to track workload")
	}
	return tasks, nil
}

// --- Lookup Table Management Implementations ---

// Staff Roles
func (uc *staffUseCaseImpl) AddStaffRole(ctx context.Context, name, description string) (*entity.StaffRole, error) {
	uc.logger.Info("Adding staff role", "name", name)
	if name == "" {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "role name cannot be empty")
	}
	role := &entity.StaffRole{Name: name, Description: description}
	err := uc.staffRoleRepo.Create(ctx, role)
	if err != nil {
		uc.logger.Error("Failed to add staff role", "name", name, "error", err)
		// TODO: Handle unique constraint violation error specifically
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to add staff role")
	}
	return role, nil
}

func (uc *staffUseCaseImpl) ListStaffRoles(ctx context.Context) ([]*entity.StaffRole, error) {
	uc.logger.Info("Listing staff roles")
	roles, err := uc.staffRoleRepo.ListAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to list staff roles", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to list staff roles")
	}
	return roles, nil
}

// Staff Statuses
func (uc *staffUseCaseImpl) AddStaffStatus(ctx context.Context, name, description string) (*entity.StaffStatus, error) {
	uc.logger.Info("Adding staff status", "name", name)
	if name == "" {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "status name cannot be empty")
	}
	status := &entity.StaffStatus{Name: name, Description: description}
	err := uc.staffStatusRepo.Create(ctx, status)
	if err != nil {
		uc.logger.Error("Failed to add staff status", "name", name, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to add staff status")
	}
	return status, nil
}

func (uc *staffUseCaseImpl) ListStaffStatuses(ctx context.Context) ([]*entity.StaffStatus, error) {
	uc.logger.Info("Listing staff statuses")
	statuses, err := uc.staffStatusRepo.ListAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to list staff statuses", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to list staff statuses")
	}
	return statuses, nil
}

// Task Statuses
func (uc *staffUseCaseImpl) AddTaskStatus(ctx context.Context, name, description string) (*entity.TaskStatus, error) {
	uc.logger.Info("Adding task status", "name", name)
	if name == "" {
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInvalidInput, "task status name cannot be empty")
	}
	status := &entity.TaskStatus{Name: name, Description: description}
	err := uc.taskStatusRepo.Create(ctx, status)
	if err != nil {
		uc.logger.Error("Failed to add task status", "name", name, "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to add task status")
	}
	return status, nil
}

func (uc *staffUseCaseImpl) ListTaskStatuses(ctx context.Context) ([]*entity.TaskStatus, error) {
	uc.logger.Info("Listing task statuses")
	statuses, err := uc.taskStatusRepo.ListAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to list task statuses", "error", err)
		return nil, coreUseCase.NewUseCaseError(coreUseCase.ErrInternal, "failed to list task statuses")
	}
	return statuses, nil
}
