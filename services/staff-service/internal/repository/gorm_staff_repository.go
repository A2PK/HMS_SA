package repository

import (
	"context"
	"errors"
	"time"

	coreRepository "golang-microservices-boilerplate/pkg/core/repository"
	"golang-microservices-boilerplate/services/staff-service/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Staff Repository ---

// Ensure GormStaffRepository implements StaffRepository
var _ StaffRepository = (*GormStaffRepository)(nil)

type GormStaffRepository struct {
	*coreRepository.GormBaseRepository[entity.Staff] // Embed the generic GORM repository
}

// NewGormStaffRepository creates a new GORM staff repository
func NewGormStaffRepository(db *gorm.DB) *GormStaffRepository {
	//nolint:all // Suppress T does not satisfy Entity constraint due to pointer receivers
	baseRepo := coreRepository.NewGormBaseRepository[entity.Staff](db)
	return &GormStaffRepository{
		GormBaseRepository: baseRepo,
	}
}

// FindByID overrides base FindByID to preload associations
func (r *GormStaffRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Staff, error) {
	staff := &entity.Staff{}
	err := r.DB.WithContext(ctx).
		Preload("Role").                 // Preload StaffRole
		Preload("Status").               // Preload StaffStatus
		Preload("Schedule").             // Preload ScheduleEntries
		Preload("Schedule.Task").        // Preload Task within ScheduleEntry
		Preload("Schedule.Task.Status"). // Preload TaskStatus within Task
		First(staff, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Use a standard error or define custom errors
			return nil, errors.New("staff not found")
		}
		return nil, err
	}
	return staff, nil
}

// FindAvailableDoctors implements the specific logic. Needs significant adjustment.
// TODO: Refactor this query based on new entity structure (RoleID, StatusID, Task Start/End times)
func (r *GormStaffRepository) FindAvailableDoctors(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entity.Staff, error) {
	var availableDoctors []*entity.Staff

	// Placeholder query - MUST BE REVISED
	// This only finds active doctors, but doesn't analyze their schedule for actual availability gaps.
	err := r.DB.WithContext(ctx).
		Preload("Role").
		Preload("Status").
		Preload("Schedule.Task.Status").
		Joins("JOIN "+entity.StaffRole{}.TableName()+" sr ON sr.name = staff.role_id").
		Joins("JOIN "+entity.StaffStatus{}.TableName()+" ss ON ss.name = staff.status_id").
		Where("sr.name = ?", "Doctor"). // Assuming "Doctor" is the role name
		Where("ss.name = ?", "Active"). // Assuming "Active" is the status name
		Where("staff.deleted_at IS NULL").
		Find(&availableDoctors).Error

	// Complex logic needed here to:
	// 1. Fetch tasks for these doctors within the timeframe [startTime, endTime]
	// 2. Calculate free time slots based on task start/end times.
	// 3. Filter the 'availableDoctors' list based on those who have free slots.

	if err != nil {
		return nil, err
	}

	// Currently returns all active doctors, NOT filtered by actual schedule availability.
	return availableDoctors, nil
}

// AddScheduleEntries creates tasks and links them via ScheduleEntry within a transaction.
func (r *GormStaffRepository) AddScheduleEntries(ctx context.Context, staffID uuid.UUID, tasks []*entity.Task) error {
	if len(tasks) == 0 {
		return nil // Nothing to add
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create all the tasks first (ensure IDs are generated)
		if err := tx.Create(tasks).Error; err != nil {
			return err // Rollback transaction
		}

		// 2. Create the ScheduleEntry links
		scheduleEntries := make([]entity.ScheduleEntry, len(tasks))
		for i, task := range tasks {
			scheduleEntries[i] = entity.ScheduleEntry{
				StaffID: staffID,
				TaskID:  task.ID, // Use the ID of the newly created task
			}
		}

		if err := tx.Create(&scheduleEntries).Error; err != nil {
			return err // Rollback transaction
		}

		return nil // Commit transaction
	})
}

// UpdateStatus updates only the StatusID for a staff member.
func (r *GormStaffRepository) UpdateStatus(ctx context.Context, staffID uuid.UUID, statusID string) error {
	result := r.DB.WithContext(ctx).
		Model(&entity.Staff{}).
		Where("id = ?", staffID).
		Updates(map[string]interface{}{
			"status_id":  statusID,
			"updated_at": time.Now(), // Ensure updated_at is set
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("staff not found or status not changed")
	}
	return nil
}

// AssignTaskToStaff creates a Task and links it via ScheduleEntry in a transaction.
func (r *GormStaffRepository) AssignTaskToStaff(ctx context.Context, staffID uuid.UUID, task *entity.Task) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create the task
		if err := tx.Create(task).Error; err != nil {
			return err
		}

		// 2. Create the ScheduleEntry link
		scheduleEntry := entity.ScheduleEntry{
			StaffID: staffID,
			TaskID:  task.ID,
		}
		if err := tx.Create(&scheduleEntry).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindTasksByStaffID retrieves tasks associated with a staff member.
func (r *GormStaffRepository) FindTasksByStaffID(ctx context.Context, staffID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.DB.WithContext(ctx).
		Preload("Status"). // Preload TaskStatus
		Joins("JOIN "+entity.ScheduleEntry{}.TableName()+" sse ON sse.task_id = "+entity.Task{}.TableName()+".id").
		Where("sse.staff_id = ?", staffID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// --- Staff Role Repository ---

// Ensure GormStaffRoleRepository implements StaffRoleRepository
var _ StaffRoleRepository = (*GormStaffRoleRepository)(nil)

type GormStaffRoleRepository struct {
	db *gorm.DB
}

func NewGormStaffRoleRepository(db *gorm.DB) *GormStaffRoleRepository {
	return &GormStaffRoleRepository{db: db}
}

func (r *GormStaffRoleRepository) Create(ctx context.Context, role *entity.StaffRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *GormStaffRoleRepository) FindByName(ctx context.Context, name string) (*entity.StaffRole, error) {
	var role entity.StaffRole
	err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff role not found") // Consistent error
		}
		return nil, err
	}
	return &role, nil
}

func (r *GormStaffRoleRepository) ListAll(ctx context.Context) ([]*entity.StaffRole, error) {
	var roles []*entity.StaffRole
	err := r.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// --- Staff Status Repository ---

// Ensure GormStaffStatusRepository implements StaffStatusRepository
var _ StaffStatusRepository = (*GormStaffStatusRepository)(nil)

type GormStaffStatusRepository struct {
	db *gorm.DB
}

func NewGormStaffStatusRepository(db *gorm.DB) *GormStaffStatusRepository {
	return &GormStaffStatusRepository{db: db}
}

func (r *GormStaffStatusRepository) Create(ctx context.Context, status *entity.StaffStatus) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *GormStaffStatusRepository) FindByName(ctx context.Context, name string) (*entity.StaffStatus, error) {
	var status entity.StaffStatus
	err := r.db.WithContext(ctx).First(&status, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff status not found")
		}
		return nil, err
	}
	return &status, nil
}

func (r *GormStaffStatusRepository) ListAll(ctx context.Context) ([]*entity.StaffStatus, error) {
	var statuses []*entity.StaffStatus
	err := r.db.WithContext(ctx).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// --- Task Status Repository ---

// Ensure GormTaskStatusRepository implements TaskStatusRepository
var _ TaskStatusRepository = (*GormTaskStatusRepository)(nil)

type GormTaskStatusRepository struct {
	db *gorm.DB
}

func NewGormTaskStatusRepository(db *gorm.DB) *GormTaskStatusRepository {
	return &GormTaskStatusRepository{db: db}
}

func (r *GormTaskStatusRepository) Create(ctx context.Context, status *entity.TaskStatus) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *GormTaskStatusRepository) FindByName(ctx context.Context, name string) (*entity.TaskStatus, error) {
	var status entity.TaskStatus
	err := r.db.WithContext(ctx).First(&status, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task status not found")
		}
		return nil, err
	}
	return &status, nil
}

func (r *GormTaskStatusRepository) ListAll(ctx context.Context) ([]*entity.TaskStatus, error) {
	var statuses []*entity.TaskStatus
	err := r.db.WithContext(ctx).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// --- Task Repository ---

// Ensure GormTaskRepository implements TaskRepository
var _ TaskRepository = (*GormTaskRepository)(nil)

type GormTaskRepository struct {
	*coreRepository.GormBaseRepository[entity.Task] // Embed the generic GORM repository
}

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{
		GormBaseRepository: coreRepository.NewGormBaseRepository[entity.Task](db),
	}
}
